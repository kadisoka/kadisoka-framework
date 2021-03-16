//

import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;

import 'rest/oauth2.dart';
import 'rest/terminal.dart';

enum IamSessionState {
  Unknown,
  SignedOut,
  SignedIn,
}

//TODO: make this a ChangeNotifier, stream source or something to make it reactive.
abstract class IamServiceClient {
  Future<String> registerTerminal(
      String userIdentifier, Iterable<String> userPreferredLanguages);

  Future<String> confirmAuthorizationVerification(String code);

  Future<bool> signOut();

  // True if we have an access token. To actually check if the session is
  // still valid, use [checkSession].
  bool get isSignedIn;

  // The identifier used to sign in.
  String get identifier;

  String get terminalId;

  String get accessToken;

  ValueNotifier<IamSessionState> get sessionStateNotifier;
}

class IamServiceClientImpl implements IamServiceClient {
  final String serverBaseUrl;
  http.Client _httpClient;

  String _clientId;
  String _clientSecret;

  String _identifier = '';
  String _terminalId = '';
  String _accessToken = '';

  IamServiceClientImpl({
    @required this.serverBaseUrl,
    http.Client httpClient,
    @required String clientId,
    @required String clientSecret,
  })  : assert(serverBaseUrl?.isNotEmpty == true),
        assert(clientId?.isNotEmpty == true),
        assert(clientSecret?.isNotEmpty == true) {
    _httpClient = httpClient ?? http.Client();
    _clientId = clientId;
    _clientSecret = clientSecret;
  }

  @override
  bool get isSignedIn =>
      _accessToken?.isNotEmpty == true &&
      _sessionStateNotifier.value == IamSessionState.SignedIn;

  @override
  String get identifier => _identifier;

  @override
  String get terminalId => _terminalId;

  @override
  String get accessToken => _accessToken;

  @override
  Future<String> registerTerminal(
    String userIdentifier,
    Iterable<String> userPreferredLanguages,
  ) async {
    if (userIdentifier?.isNotEmpty != true) {
      return null;
    }

    _identifier = userIdentifier;

    final resp = await _httpClient.post(
      Uri.http(serverBaseUrl, 'rest/v1/terminals/register'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': _httpClientAuthorization,
        'Accept-Language': userPreferredLanguages.join(', '),
      },
      body: jsonEncode(TerminalRegisterPostRequestJsonV1(
        verificationResourceName: userIdentifier,
        verificationMethods: <String>['none'],
      ).toRestV1Json()),
    );
    if (resp.statusCode != 200) {
      throw Exception('Fetch got ${resp.statusCode}');
    }

    final respData = TerminalRegisterPostResponseJsonV1.fromRestV1Json(
        jsonDecode(resp.body));
    _terminalId = respData.terminalId;

    return _terminalId;
  }

  @override
  Future<String> confirmAuthorizationVerification(String code) async {
    final resp = await _httpClient.post(
      Uri.http(serverBaseUrl, 'rest/v1/oauth2/token'),
      headers: <String, String>{
        'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
        'Authorization': _httpClientAuthorization,
      },
      body: <String, dynamic>{
        'grant_type': 'authorization_code',
        'code': 'otp:$terminalId:$code',
      },
    );
    if (resp.statusCode != 200) {
      throw Exception('Fetch got ${resp.statusCode}');
    }

    final respData = OAuth2TokenResponse.fromJson(jsonDecode(resp.body));
    _accessToken = respData.accessToken;

    _sessionStateNotifier.value = IamSessionState.SignedIn;

    return _accessToken;
  }

  @override
  Future<bool> signOut() async {
    final resp = await _httpClient.delete(
      Uri.http(serverBaseUrl, 'rest/v1/terminals/self'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': _httpBearerAuthorization,
      },
      body: '{}',
    );
    if (resp.statusCode != 200) {
      throw Exception('Fetch got ${resp.statusCode}');
    }

    _accessToken = '';
    _terminalId = '';

    _sessionStateNotifier.value = IamSessionState.SignedOut;

    return true;
  }

  String get _httpClientAuthorization =>
      'Basic ' + base64Encode(utf8.encode('$_clientId:$_clientSecret'));

  String get _httpBearerAuthorization => 'Bearer ' + _accessToken;

  final ValueNotifier<IamSessionState> _sessionStateNotifier =
      ValueNotifier(IamSessionState.Unknown);

  @override
  ValueNotifier<IamSessionState> get sessionStateNotifier =>
      _sessionStateNotifier;
}
