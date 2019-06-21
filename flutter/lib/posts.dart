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
  State<StatefulWidget> createState() => postsState();
}

class PostDataSource extends DataTableSource {
  List<Map<String, dynamic>> posts;
  final Function switchToPostPage;
  final String postType;
  String searchTerm;
  int rowsPerPage = 1;
  int pageNum = 0;
  bool sortAscending = true;
  int sortColumnIndex = 0;
  List<int> availableRowsPerPage = new List<int>();
  PostDataSource({@required this.switchToPostPage, @required this.postType});

  void _setAvailableRowsPerPage(List<int> newRowsPerPage) {
    availableRowsPerPage = newRowsPerPage;
  }

  void _goToPostPage(String postid) {
    print('navigate to page $postid with posttype $postType');
    switchToPostPage(postid);
  }

  Future<List<Map<String, dynamic>>> _getPosts() async {
    print('start query');
    String searchTermQueryString = "";
    if (searchterm != null) {
      searchTermQueryString = 'searchterm:"$searchTerm",';
    }
    Map<String, String> queryParameters = {
      'query':
          '{posts(type:"$postType",perpage:$rowsPerPage,page:$pageNum,${searchTermQueryString}sort:"title",ascending:$sortAscending){title views id author date}}',
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

  @override
  DataRow getRow(int index) {
    Map<String, dynamic> post = posts[index];
    print('id ${post['id']}');
    return DataRow.byIndex(index: index, cells: <DataCell>[
      DataCell(Text('title ${post['title']}'), onTap: () {
        print('tapped');
        _goToPostPage(post['id']);
      }),
      DataCell(Text('date ${post['date']}'),
          onTap: () => _goToPostPage(post['id'])),
      DataCell(Text('views ${post['views']}'),
          onTap: () => _goToPostPage(post['id']))
    ]);
  }

  @override
  bool get isRowCountApproximate => false;

  @override
  int get rowCount => posts.length;

  void sort(int columnIndex, bool ascending) {
    print('sort column $columnIndex');
    notifyListeners();
  }

  void setPage(int newPageNum) {
    pageNum = newPageNum;
    print('switch to page $pageNum');
    notifyListeners();
  }

  void setRowsPerPage(int newRowsPerPage) {
    rowsPerPage = newRowsPerPage;
    print('rows per page $rowsPerPage');
    notifyListeners();
  }

  @override
  int get selectedRowCount => 0;
}

class postsState extends State<PostsPage> {
  PostDataSource postDataSource;
  String postType;
  final GlobalKey<FormState> formKey = new GlobalKey<FormState>();

  @override
  void initState() {
    super.initState();
    postType = widget.postType;
    print('got post type $postType');
    postDataSource = PostDataSource(
      postType: postType,
      switchToPostPage: widget.switchToPostPage,
    );
    postDataSource._setAvailableRowsPerPage([postDataSource.rowsPerPage]);
    postDataSource._getPosts().then((res) {
      setState(() {
        postDataSource.posts = res;
      });
    }).catchError((err) {
      print(err);
    });
  }

  @override
  Widget build(BuildContext context) {
    if (postDataSource.posts != null) {
      return Container(
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextFormField(
                validator: (value) {
                  if (value.isEmpty) {
                    return 'enter search query';
                  }
                  return null;
                },
              ),
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 16.0),
                child: RaisedButton(
                  onPressed: () {
                    if (_formKey.currentState.validate()) {
                      Scaffold.of(context)
                          .showSnackBar(SnackBar(content: Text('searching...')));
                      setState(() {
                        postDataSource.searchTerm = _formKey.currentState.value;
                      });
                      postDataSource._getPosts().then((res) {
                        setState(() {
                          postDataSource.posts = res;
                        });
                      }).catchError((err) {
                        print(err);
                      });
                    }
                  },
                  child: Text('Submit'),
                ),
              ),
              PaginatedDataTable(
                  header: Text('${postType}s'),
                  availableRowsPerPage: postDataSource.availableRowsPerPage,
                  rowsPerPage: postDataSource.rowsPerPage,
                  onRowsPerPageChanged: postDataSource.setRowsPerPage,
                  sortColumnIndex: postDataSource.sortColumnIndex,
                  sortAscending: postDataSource.sortAscending,
                  onPageChanged: postDataSource.setPage,
                  columns: <DataColumn>[
                    DataColumn(
                        label: const Text('Title'),
                        onSort: postDataSource.sort),
                    DataColumn(
                        label: const Text('Date'),
                        onSort: postDataSource.sort),
                    DataColumn(
                        label: const Text('Views'),
                        onSort: postDataSource.sort)
                  ],
                  source: postDataSource)
            ],
          ),
        ),
      );
    } else {
      return CircularProgressIndicator();
    }
  }
}
