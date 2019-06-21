import 'package:flutter/material.dart';
import 'posts.dart';
import 'post.dart';
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
  String postid;
  bool viewPost = false;

  void switchToPostPage(String newpostid) {
    setState(() {
      postid = newpostid;
      viewPost = true;
    });
  }

  @override
  Widget build(BuildContext context) {
    Widget bodyWidget;
    switch (index) {
      case 0:
        bodyWidget = HomePage();
        break;
      case 1:
        if (!viewPost) {
          bodyWidget = BlogsPage(switchToPostPage: switchToPostPage);
        } else {
          bodyWidget = BlogPage(blogid: postid);
        }
        break;
      case 2:
        if (!viewPost) {
          bodyWidget = ProjectsPage(switchToPostPage: switchToPostPage);
        } else {
          bodyWidget = ProjectPage(projectid: postid);
        }
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
              () {
                index = newIndex;
                viewPost = false;
                print('new index $index');
              },
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
