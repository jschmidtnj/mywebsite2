import 'package:flutter/material.dart';

class HomePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      child: Center(
        child: Text(
          'Joshua',
          style: Theme.of(context).textTheme.display2,
        ),
      ),
    );
  }
}
