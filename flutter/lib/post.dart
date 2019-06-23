import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'config.dart';
import 'dart:io';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:intl/intl.dart';

class ProjectPage extends PostPage {
  final String projectid;
  ProjectPage({Key key, @required this.projectid})
      : super(key: key, postid: projectid, postType: 'project');
}

class BlogPage extends PostPage {
  final String blogid;
  BlogPage({Key key, @required this.blogid})
      : super(key: key, postid: blogid, postType: 'blog');
}

// needs scaffold from main nav class - but no circular dependency
class PostPage extends StatelessWidget {
  final String postid;
  final String postType;

  PostPage({Key key, @required this.postid, @required this.postType})
      : super(key: key);

  Future<Map<String, dynamic>> _getPostData() async {
    var queryParameters = {
      'query':
          '{post(type:"$postType",id:"$postid"){title views id author date content}}'
    };
    var uri = Uri.https(config['apiURL'], '/graphql', queryParameters);
    final response = await http.get(uri, headers: {
      HttpHeaders.contentTypeHeader: 'application/json',
    });
    print(response.body);
    if (response.statusCode == 200) {
      return json.decode(response.body);
    } else {
      throw Exception('Failed to load post');
    }
  }

  @override
  StatelessWidget build(BuildContext context) {
    print('got id $postid');
    final Future<Map<String, dynamic>> postdata = _getPostData();
    return Container(
      child: FutureBuilder<Map<String, dynamic>>(
        future: postdata,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            RegExp imgRegex = new RegExp("<img([\\w\\W]+?)><\\/img>",
                caseSensitive: true, multiLine: true);
            RegExp srcRegex = new RegExp("src\\s*=\\s*\"(.+?)\"",
                caseSensitive: true, multiLine: true);
            String content =
                Uri.decodeComponent(snapshot.data['data']['post']['content']);
            for (Match imgMatch in imgRegex.allMatches(content)) {
              print(
                  'got image match start ${imgMatch.start} end ${imgMatch.end}');
              String imgTag = content.substring(imgMatch.start, imgMatch.end);
              print('got tag $imgTag');
              Match srcMatch = srcRegex.firstMatch(imgTag);
              String imgSrcWithQuotes = imgTag
                  .substring(srcMatch.start, srcMatch.end)
                  .split('src=')
                  .last;
              String imgSrc =
                  imgSrcWithQuotes.substring(1, imgSrcWithQuotes.length - 1);
              print('got src $imgSrc');
              String imgName = imgSrc.split('/').last.split('.').first;
              print('got image name $imgName');
              content = content.replaceFirst(imgTag, '![$imgName]($imgSrc)');
              print('new content $content');
            }
            String title =
                Uri.decodeComponent(snapshot.data['data']['post']['title']);
            String author =
                Uri.decodeComponent(snapshot.data['data']['post']['author']);
            int views = snapshot.data['data']['post']['views'];
            DateTime date = new DateTime.fromMillisecondsSinceEpoch(int.parse(
                    snapshot.data['data']['post']['id']
                        .toString()
                        .substring(0, 8),
                    radix: 16) *
                1000);
            DateFormat formatter = new DateFormat('yyyy-MM-dd');
            String dateStr = formatter.format(date);
            return Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Padding(padding: EdgeInsets.only(top: 40.0)),
                  Text(title,
                      style: TextStyle(
                        fontSize: 25.0,
                        fontWeight: FontWeight.bold,
                      )),
                  Padding(padding: EdgeInsets.only(top: 10.0)),
                  Text('by $author. $dateStr. $views views'),
                  Padding(padding: EdgeInsets.only(top: 10.0)),
                  Flexible(child: Markdown(data: content))
                ]);
          } else if (snapshot.hasError) {
            return Text("${snapshot.error}");
          }
          // By default, show a loading spinner.
          return CircularProgressIndicator();
        },
      ),
    );
  }
}

/// The first line of a HTTP request body
class Quotation extends StatefulWidget {
  Quotation({Key key, this.url}) : super(key: key);

  final String url;

  @override
  createState() => _QuotationState();
}

class _QuotationState extends State<Quotation> {
  String data = 'Loading ...';

  @override
  void initState() {
    super.initState();
    _get();
  }

  _get() async {
    final res = await http.get(widget.url);
    setState(() => data = _parseQuoteFromJson(res.body));
  }

  String _parseQuoteFromJson(String jsonStr) {
    // In the real world, this should check for errors
    final jsonQuote = json.decode(jsonStr);
    return jsonQuote['contents']['quotes'][0]['quote'];
  }

  @override
  Widget build(BuildContext context) {
    return Text(data, textAlign: TextAlign.center);
  }
}