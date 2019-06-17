import 'package:flutter/material.dart';

class ProjectsPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      child: Center(
        child: Text(
          'Projects',
          style: Theme.of(context).textTheme.display2,
        ),
      ),
    );
  }
}
