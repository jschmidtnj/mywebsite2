import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';

class Loading extends StatelessWidget {
  @override
  StatelessWidget build(BuildContext context) {
    return Container(child: Center(child: CircularProgressIndicator()));
  }
}
