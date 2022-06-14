//

import 'package:flutter/foundation.dart';

class TerminalRegisterPostRequestJsonV1 {
  final String verificationResourceName;
  final Iterable<String> verificationMethods;
  final String displayName;

  TerminalRegisterPostRequestJsonV1({
    @required this.verificationResourceName,
    this.verificationMethods,
    this.displayName,
  }) : assert(verificationResourceName?.isNotEmpty == true);

  Map<String, dynamic> toRestV1Json() {
    return {
      'verification_resource_name': verificationResourceName,
      'verification_methods': verificationMethods,
      'display_name': displayName,
    };
  }
}

class TerminalRegisterPostResponseJsonV1 {
  final String terminalId;
  final String terminalSecret;
  final DateTime codeExpiry;

  TerminalRegisterPostResponseJsonV1({
    @required this.terminalId,
    this.terminalSecret,
    this.codeExpiry,
  }) : assert(terminalId?.isNotEmpty == true);

  factory TerminalRegisterPostResponseJsonV1.fromRestV1Json(
      Map<String, dynamic> json) {
    return TerminalRegisterPostResponseJsonV1(
      terminalId: json['terminal_id'],
      terminalSecret: json['terminal_secret'],
      codeExpiry: DateTime.parse(json['code_expiry']),
    );
  }
}
