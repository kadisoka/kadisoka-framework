//

import 'package:flutter/material.dart';

import 'kadisoka/iam/service.dart';
import 'sign_in_otp_page.dart';

class SignInIdentifierPage extends StatefulWidget {
  static MaterialPageRoute createRoute({
    RouteSettings routeSettings,
    bool fullscreenDialog = false,
    Key widgetKey,
    IamServiceClient iamService,
  }) =>
      MaterialPageRoute(
          settings: routeSettings,
          fullscreenDialog: fullscreenDialog,
          builder: (BuildContext context) => SignInIdentifierPage(
                key: widgetKey,
                iamService: iamService,
                fullscreenDialog: fullscreenDialog,
              ));

  SignInIdentifierPage({
    Key key,
    @required this.iamService,
    this.fullscreenDialog,
  })  : assert(iamService != null),
        super(key: key);

  final IamServiceClient iamService;
  final bool fullscreenDialog;

  @override
  _SignInIdentifierPageState createState() => _SignInIdentifierPageState();
}

class _SignInIdentifierPageState extends State<SignInIdentifierPage> {
  final _identifierTextController = TextEditingController();

  @override
  void initState() {
    super.initState();

    WidgetsBinding.instance.addPostFrameCallback((_) {
      final nav = Navigator.of(context);
      if (widget.iamService?.accessToken?.isNotEmpty == true) {
        //TODO: check with the service if the access token is still usable.
        while (nav.canPop()) {
          nav.pop();
        }
        return;
      }
      if (widget.iamService?.terminalId?.isNotEmpty == true) {
        nav.pushReplacement(SignInOtpPage.createRoute(
          fullscreenDialog: widget.fullscreenDialog,
          iamService: widget.iamService,
        ));
      }
    });
  }

  @override
  void dispose() {
    _identifierTextController?.dispose();

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
      title: Text('Sign In'),
    );
  }

  Widget _buildBody(BuildContext context) {
    const rowHPadding = 16.0;
    const vSpacer = SizedBox(height: 8);
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
                controller: _identifierTextController,
                decoration: InputDecoration(
                  labelText: 'Phone or email',
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
                    onPressed: () {
                      final identifier = _identifierTextController.text;
                      if (identifier.isEmpty) {
                        return;
                      }

                      widget.iamService
                          .registerTerminal(identifier)
                          .then((String terminalId) {
                        if (terminalId?.isNotEmpty == true) {
                          Navigator.of(context).pushReplacement(
                            SignInOtpPage.createRoute(
                              fullscreenDialog: widget.fullscreenDialog,
                              iamService: widget.iamService,
                            ),
                          );
                        }
                      });
                    },
                    child: Text('Next'),
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
