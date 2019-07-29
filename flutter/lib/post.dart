import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'config.dart';
import 'dart:io';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:intl/intl.dart';
import 'package:html/parser.dart' show parse;
import 'package:html/dom.dart' as dom;

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
          '{post(type:"$postType",id:"$postid",cache:true){title views id author date content}}'
    };
    var uri = Uri.https(config['apiURL'], '/graphql', queryParameters);
    final response = await http.get(uri, headers: {
      HttpHeaders.contentTypeHeader: 'application/json',
    });
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
            String content =
                Uri.decodeComponent(snapshot.data['data']['post']['content'])
                    .replaceAll('</img>', '');
            dom.Document document = parse(content);
            List<dom.Element> linkTags = document.querySelectorAll('img.lazy');
            for (dom.Element linkTag in linkTags) {
              String tagToReplace = linkTag.outerHtml;
              String imgSrc = linkTag.attributes["data-src"];
              String alt = linkTag.attributes["alt"];
              content = content.replaceFirst(tagToReplace, '![$alt]($imgSrc)');
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
