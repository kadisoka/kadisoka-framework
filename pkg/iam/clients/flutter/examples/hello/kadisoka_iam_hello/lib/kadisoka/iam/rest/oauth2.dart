//

import 'package:flutter/foundation.dart';

class OAuth2TokenResponse {
  final String accessToken;

  OAuth2TokenResponse({
    @required this.accessToken,
  }) : assert(accessToken?.isNotEmpty == true);

  factory OAuth2TokenResponse.fromJson(Map<String, dynamic> json) {
    return OAuth2TokenResponse(
      accessToken: json['access_token'],
    );
  }
}
