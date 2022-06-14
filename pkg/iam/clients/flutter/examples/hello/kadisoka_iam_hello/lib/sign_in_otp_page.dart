//

import 'package:flutter/material.dart';

import 'kadisoka/iam/service.dart';

class SignInOtpPage extends StatefulWidget {
  static MaterialPageRoute createRoute({
    RouteSettings routeSettings,
    bool fullscreenDialog = false,
    Key widgetKey,
    IamServiceClient iamService,
  }) =>
      MaterialPageRoute(
          settings: routeSettings,
          fullscreenDialog: fullscreenDialog,
          builder: (BuildContext context) => SignInOtpPage(
                key: widgetKey,
                iamService: iamService,
              ));

  SignInOtpPage({
    Key key,
    @required this.iamService,
  })  : assert(iamService != null),
        super(key: key);

  final IamServiceClient iamService;

  @override
  _SignInOtpPageState createState() => _SignInOtpPageState();
}

class _SignInOtpPageState extends State<SignInOtpPage> {
  final _otpTextController = TextEditingController();

  @override
  void initState() {
    super.initState();

    if (widget.iamService?.identifier?.isNotEmpty == true) {}
  }

  @override
  void dispose() {
    _otpTextController?.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(context),
      body: _buildBody(context),
    );
  }

  Widget _buildAppBar(BuildContext context) {
    return AppBar(
      title: Text('Sign In — Confirm Verification'),
    );
  }

  Widget _buildBody(BuildContext context) {
    const rowHPadding = 16.0;
    const vSpacer = SizedBox(height: 8);
    const hSpacer = SizedBox(width: 4);
    return Center(
      child: ConstrainedBox(
        constraints: BoxConstraints(maxWidth: 440),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: rowHPadding),
              child: TextFormField(
                controller: _otpTextController,
                decoration: InputDecoration(
                  labelText: 'Verification code',
                  border: OutlineInputBorder(),
                ),
              ),
            ),
            vSpacer,
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: rowHPadding),
              child: Row(
                children: [
                  Spacer(),
                  ElevatedButton(
                      onPressed: () {},
                      child: Text('Use different email/phone')),
                  hSpacer,
                  ElevatedButton(
                    onPressed: () {
                      final otp = _otpTextController.text;
                      if (otp.isEmpty) {
                        return;
                      }

                      widget.iamService
                          .confirmAuthorizationVerification(otp)
                          .then((accessToken) {
                        final nav = Navigator.of(context);
                        while (nav.canPop()) {
                          nav.pop();
                        }
                      });
                    },
                    child: Text('Confirm'),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
