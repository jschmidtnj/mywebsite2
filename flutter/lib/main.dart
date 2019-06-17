import 'package:flutter/material.dart';
import 'projects.dart';
import 'blogs.dart';
import 'homepage.dart';

void main() => runApp(MyApp());

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'Flutter Bottom Nav Demo',
        theme: ThemeData(
          primarySwatch: Colors.deepOrange,
        ),
        home: Nav());
  }
}

class Nav extends StatefulWidget {
  @override
  createState() => NavState();
}

class NavState extends State<Nav> {
  int index = 0;

  @override
  Widget build(BuildContext context) {
    StatelessWidget bodyWidget;
    switch (this.index) {
      case 0:
        bodyWidget = HomePage();
        break;
      case 1:
        bodyWidget = BlogsPage();
        break;
      case 2:
        bodyWidget = ProjectsPage();
        break;
      default:
        bodyWidget = HomePage();
        break;
    }
    return Scaffold(
      body: bodyWidget,
      bottomNavigationBar: BottomNavBar(
        index: index,
        callback: (newIndex) => setState(
              () => this.index = newIndex,
            ),
      ),
    );
  }
}

class Body extends StatelessWidget {
  Body(this.index);
  final int index;
  @override
  Widget build(BuildContext context) {
    return Container(
      child: Center(
        child: Text(
          'Index $index',
          style: Theme.of(context).textTheme.display2,
        ),
      ),
    );
  }
}

class BottomNavBar extends StatelessWidget {
  BottomNavBar({this.index, this.callback});
  final int index;
  final Function(int) callback;

  @override
  Widget build(BuildContext context) {
    /// BottomNavigationBar is automatically set to type 'fixed'
    /// when there are three of less items
    return BottomNavigationBar(
      currentIndex: index,
      onTap: callback,
      items: <BottomNavigationBarItem>[
        BottomNavigationBarItem(
          icon: Icon(Icons.home),
          title: Text('Home'),
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.search),
          title: Text('Blogs'),
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.search),
          title: Text('Projects'),
        ),
      ],
    );
  }
}
