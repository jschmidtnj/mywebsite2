import 'package:flutter/material.dart';
import 'posts.dart';
import 'post.dart';

/* main flutter class entry point for app. */
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

  Future<bool> goBack() {
    return new Future<bool>(() {
      setState(() {
        viewPost = false;
      });
      return false;
    });
  }

  @override
  Widget build(BuildContext context) {
    Widget bodyWidget;
    switch (index) {
      case 0:
        if (!viewPost) {
          bodyWidget = BlogsPage(switchToPostPage: switchToPostPage);
        } else {
          bodyWidget =
              WillPopScope(onWillPop: goBack, child: BlogPage(blogid: postid));
        }
        break;
      case 1:
        if (!viewPost) {
          bodyWidget = ProjectsPage(switchToPostPage: switchToPostPage);
        } else {
          bodyWidget = WillPopScope(
              onWillPop: goBack, child: ProjectPage(projectid: postid));
        }
        break;
      default:
        bodyWidget = BlogsPage(switchToPostPage: switchToPostPage);
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
          icon: Icon(Icons.edit),
          title: Text('Blogs'),
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.edit_attributes),
          title: Text('Projects'),
        ),
      ],
    );
  }
}
