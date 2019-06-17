import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'config.dart';

Future<Map<String, dynamic>> _getBlogData() async {
  final response = await http.get(config['apiURL'] + '/hello');
  if (response.statusCode == 200) {
    return json.decode(response.body);
  } else {
    throw Exception('Failed to load post');
  }
}

class BlogsPage extends StatelessWidget {
  final Future<Map<String, dynamic>> blogdata = _getBlogData();

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Padding(padding: EdgeInsets.only(top: 50.0)),
            Text(
              'Blogs',
              style: Theme.of(context).textTheme.display2,
            ),
            FutureBuilder<Map<String, dynamic>>(
              future: blogdata,
              builder: (context, snapshot) {
                if (snapshot.hasData) {
                  return Text(snapshot.data['message']);
                } else if (snapshot.hasError) {
                  return Text("${snapshot.error}");
                }
                // By default, show a loading spinner.
                return CircularProgressIndicator();
              },
            ),
          ],
        ),
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
