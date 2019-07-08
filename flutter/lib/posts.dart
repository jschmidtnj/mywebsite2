import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'config.dart';
import 'dart:io';

class ProjectsPage extends PostsPage {
  final Function switchToPostPage;
  ProjectsPage({Key key, @required this.switchToPostPage})
      : super(
            key: key, postType: 'project', switchToPostPage: switchToPostPage);
}

class BlogsPage extends PostsPage {
  final Function switchToPostPage;
  BlogsPage({Key key, @required this.switchToPostPage})
      : super(key: key, postType: 'blog', switchToPostPage: switchToPostPage);
}

class PostsPage extends StatefulWidget {
  final String postType;
  final Function switchToPostPage;
  PostsPage({Key key, @required this.postType, @required this.switchToPostPage})
      : super(key: key);

  @override
  State<StatefulWidget> createState() => PostsState();
}

class PostsState extends State<PostsPage> {
  String postType;
  int count = 0;
  final GlobalKey<FormState> formKey = GlobalKey<FormState>();
  List<Map<String, dynamic>> posts = [];
  List<Widget> footerWidgets = [];
  bool loading = true;
  String searchTerm;
  int pageNum = 0;
  bool sortAscending = true;
  int sortColumnIndex = 0;
  List<int> availableRowsPerPage = [5, 10, 15];
  int rowsPerPage;
  final List<String> columns = ['title', 'date', 'views'];

  @override
  void initState() {
    super.initState();
    postType = widget.postType;
    print('got post type $postType');
    rowsPerPage = availableRowsPerPage[0];
    _getPosts().then((res) {
      _getRowCount().then((newcount) {
        setState(() {
          count = newcount;
          posts = res;
          loading = false;
        });
      }).catchError((err1) {
        print(err1);
      });
    }).catchError((err) {
      print(err);
    });
  }

  void _goToPostPage(String postid) {
    print('navigate to page $postid');
    widget.switchToPostPage(postid);
  }

  Future<List<Map<String, dynamic>>> _getPosts(
      {int thePageNum, int theRowsPerPage}) async {
    if (thePageNum == null) {
      thePageNum = pageNum;
    }
    if (theRowsPerPage == null) {
      theRowsPerPage = rowsPerPage;
    }
    print('start query');
    String searchTermQueryString = "";
    if (searchTerm != null) {
      searchTermQueryString = 'searchterm:"$searchTerm",';
    }
    List<String> tags = new List<String>();
    List<String> categories = new List<String>();
    Map<String, String> queryParameters = {
      'query':
          '{posts(type:"$postType",perpage:$theRowsPerPage,page:$thePageNum,${searchTermQueryString}sort:"${columns[sortColumnIndex]}",ascending:$sortAscending,tags:[${tags.join(',')}],categories:[${categories.join(',')}],cache:true){title views id author date}}',
    };
    Uri uri = Uri.https(config['apiURL'], '/graphql', queryParameters);
    print('get response');
    final response = await http.get(uri, headers: {
      HttpHeaders.contentTypeHeader: 'application/json',
    });
    if (response.statusCode == 200) {
      print(response.body);
      Map<String, dynamic> resdata = json.decode(response.body);
      if (resdata['data'] != null && resdata['data']['posts'] != null) {
        return (resdata['data']['posts'] as List).cast<Map<String, dynamic>>();
      } else {
        throw Exception('data and posts not found');
      }
    } else {
      print('error');
      throw Exception('Failed to load post');
    }
  }

  void sort(int columnIndex, bool ascending) {
    sortColumnIndex = columnIndex;
    print('sort column $sortColumnIndex');
    sortAscending = ascending;
    _getPosts().then((res) {
      setState(() {
        posts = res;
        loading = false;
      });
    }).catchError((err) {
      print(err);
    });
  }

  void setRowsPerPage(int newRowsPerPage) {
    print('rows per page $newRowsPerPage');
    _getPosts(theRowsPerPage: newRowsPerPage).then((res) {
      _getRowCount().then((newcount) {
        setState(() {
          rowsPerPage = newRowsPerPage;
          posts = res;
          count = newcount;
          loading = false;
        });
      }).catchError((err1) {
        print(err1);
      });
    }).catchError((err) {
      print(err);
    });
  }

  Future<int> _getRowCount() async {
    print('start query');
    String querySearchTerm = '';
    if (searchTerm != null) {
      querySearchTerm = searchTerm;
    }
    List<String> tags = const [];
    String tagsStr = tags.join(',tags=');
    List<String> categories = const [];
    String categoriesStr = categories.join(',categories=');
    Map<String, String> queryParameters = {
      'type': postType,
      'searchterm': querySearchTerm,
      'tags': tagsStr,
      'categories': categoriesStr
    };
    Uri uri = Uri.https(config['apiURL'], '/countPosts', queryParameters);
    print('get response');
    final response = await http.get(uri, headers: {
      HttpHeaders.contentTypeHeader: 'application/json',
      HttpHeaders.acceptEncodingHeader: 'application/json'
    });
    if (response.statusCode == 200) {
      print(response.body);
      Map<String, dynamic> resdata = json.decode(response.body);
      if (resdata['count'] != null) {
        return resdata['count'];
      } else {
        throw Exception('count not found');
      }
    } else {
      print('error');
      throw Exception('Failed to load post');
    }
  }

  void setPage(int newPageNum) {
    print('switch to page $newPageNum');
    _getPosts(thePageNum: newPageNum).then((res) {
      List<Map<String, dynamic>> newPosts = res;
      print('len newposts ${newPosts.length}');
      _getRowCount().then((newcount) {
        print('got count $newcount');
        setState(() {
          count = newcount;
          posts = newPosts;
          pageNum = newPageNum;
          loading = false;
        });
      }).catchError((err1) {
        print(err1);
      });
    }).catchError((err) {
      print(err);
    });
  }

  void handlePrevious() {
    setPage(pageNum - 1);
  }

  void handleNext() {
    setPage(pageNum + 1);
  }

  @override
  Widget build(BuildContext context) {
    Widget datatable;
    if (!loading) {
      if (posts.isEmpty) {
        datatable = Text('no ${postType}s found');
      } else {
        datatable = DataTable(
            columns: <DataColumn>[
              DataColumn(label: const Text('Title'), onSort: sort),
              DataColumn(label: const Text('Date'), onSort: sort),
              DataColumn(label: const Text('Views'), onSort: sort)
            ],
            rows: posts
                .map(
                  (post) => DataRow(
                        cells: [
                          DataCell(Text(post['title']),
                              onTap: () => _goToPostPage(post['id'])),
                          DataCell(Text(post['date']),
                              onTap: () => _goToPostPage(post['id'])),
                          DataCell(Text('${post['views']}'),
                              onTap: () => _goToPostPage(post['id']))
                        ],
                      ),
                )
                .toList(),
            sortColumnIndex: sortColumnIndex,
            sortAscending: sortAscending);
        final List<Widget> availableRowsPerPageWidget =
            availableRowsPerPage.map<DropdownMenuItem<int>>((int option) {
          return DropdownMenuItem<int>(
            value: option,
            child: Text('$option'),
          );
        }).toList();
        footerWidgets = [
          Container(
              width:
                  14.0), // to match trailing padding in case we overflow and end up scrolling
          Text('num per page'),
          ConstrainedBox(
            constraints: const BoxConstraints(
                minWidth: 64.0), // 40.0 for the text, 24.0 for the icon
            child: Align(
              alignment: AlignmentDirectional.centerEnd,
              child: DropdownButtonHideUnderline(
                child: DropdownButton<int>(
                  items: availableRowsPerPageWidget,
                  value: rowsPerPage,
                  onChanged: setRowsPerPage,
                  iconSize: 24.0,
                ),
              ),
            ),
          ),
          Container(width: 32.0),
          Text('page ${pageNum + 1}'),
          Container(width: 32.0),
          IconButton(
            icon: const Icon(Icons.chevron_left),
            padding: EdgeInsets.zero,
            onPressed: pageNum <= 0 ? null : handlePrevious,
          ),
          Container(width: 24.0),
          IconButton(
            icon: const Icon(Icons.chevron_right),
            padding: EdgeInsets.zero,
            onPressed: (rowsPerPage == posts.length &&
                    count > rowsPerPage * (pageNum + 1))
                ? handleNext
                : null,
          ),
          Container(width: 14.0),
        ];
      }
    } else {
      datatable = CircularProgressIndicator();
    }
    return Container(
      child: SingleChildScrollView(
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Padding(padding: EdgeInsets.only(top: 40.0)),
              Form(
                key: formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: <Widget>[
                    TextFormField(
                      validator: (value) {
                        if (value.isEmpty) {
                          return 'enter search query';
                        }
                        return null;
                      },
                      onSaved: (String val) => searchTerm = val,
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(vertical: 16.0),
                      child: RaisedButton(
                        onPressed: () {
                          if (formKey.currentState.validate()) {
                            formKey.currentState.save();
                            Scaffold.of(context).showSnackBar(SnackBar(
                                content: Text('search for $searchTerm')));
                            setState(() {
                              loading = true;
                            });
                            _getPosts().then((res) {
                              setState(() {
                                posts = res;
                                loading = false;
                              });
                            }).catchError((err) {
                              print(err);
                            });
                          }
                        },
                        child: Text('Submit'),
                      ),
                    ),
                  ],
                ),
              ),
              datatable,
              IconTheme.merge(
                data: const IconThemeData(opacity: 0.54),
                child: Container(
                  height: 56.0,
                  child: SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    reverse: true,
                    child: Row(
                      children: footerWidgets,
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
