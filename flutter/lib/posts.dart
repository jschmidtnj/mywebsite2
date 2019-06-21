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
  State<StatefulWidget> createState() => _PostsState();
}

class PostDataSource extends DataTableSource {
  List<Map<String, dynamic>> _posts;
  final Function switchToPostPage;
  final String postType;
  int _rowsPerPage = 1;
  int _pageNum = 0;
  bool _sortAscending = true;
  int _sortColumnIndex = 0;
  List<int> _availableRowsPerPage = new List<int>();
  PostDataSource({@required this.switchToPostPage, @required this.postType});

  void _setAvailableRowsPerPage(List<int> newRowsPerPage) {
    _availableRowsPerPage = newRowsPerPage;
  }

  void _goToPostPage(String postid) {
    print('navigate to page $postid with posttype $postType');
    switchToPostPage(postid);
  }

  Future<List<Map<String, dynamic>>> _getPosts() async {
    print('start query');
    Map<String, String> queryParameters = {
      'query':
          '{posts(type:"$postType",perpage:10,page:0,searchterm:"asdf",sort:"title",ascending:false){title views id author date}}',
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
    Map<String, dynamic> post = _posts[index];
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
  int get rowCount => _posts.length;

  void _sort(int columnIndex, bool ascending) {
    print('sort column $columnIndex');
    notifyListeners();
  }

  void _setPage(int newPageNum) {
    _pageNum = newPageNum;
    print('switch to page $_pageNum');
    notifyListeners();
  }

  void _setRowsPerPage(int newRowsPerPage) {
    _rowsPerPage = newRowsPerPage;
    print('rows per page $_rowsPerPage');
    notifyListeners();
  }

  @override
  int get selectedRowCount => 0;
}

class _PostsState extends State<PostsPage> {
  PostDataSource _postDataSource;
  String postType;

  @override
  void initState() {
    super.initState();
    postType = widget.postType;
    print('got post type $postType');
    _postDataSource = PostDataSource(
      postType: postType,
      switchToPostPage: widget.switchToPostPage,
    );
    _postDataSource._setAvailableRowsPerPage([_postDataSource._rowsPerPage]);
    _postDataSource._getPosts().then((res) {
      setState(() {
        _postDataSource._posts = res;
      });
    }).catchError((err) {
      print(err);
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_postDataSource._posts != null) {
      return Container(
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              PaginatedDataTable(
                  header: Text('${postType}s'),
                  availableRowsPerPage: _postDataSource._availableRowsPerPage,
                  rowsPerPage: _postDataSource._rowsPerPage,
                  onRowsPerPageChanged: _postDataSource._setRowsPerPage,
                  sortColumnIndex: _postDataSource._sortColumnIndex,
                  sortAscending: _postDataSource._sortAscending,
                  onPageChanged: _postDataSource._setPage,
                  columns: <DataColumn>[
                    DataColumn(
                        label: const Text('Title'),
                        onSort: _postDataSource._sort),
                    DataColumn(
                        label: const Text('Date'),
                        onSort: _postDataSource._sort),
                    DataColumn(
                        label: const Text('Views'),
                        onSort: _postDataSource._sort)
                  ],
                  source: _postDataSource)
            ],
          ),
        ),
      );
    } else {
      return CircularProgressIndicator();
    }
  }
}
