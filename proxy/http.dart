import 'utils.dart';
import 'package:http/http.dart' as http;
import 'package:shelf/shelf.dart' as shelf;
import 'package:pedantic/pedantic.dart';

Future<Content> fetchHttp(String uri, shelf.Request request) async {
  var httpClientRequest = http.StreamedRequest('GET', Uri.parse(uri));

  unawaited(request
      .read()
      .forEach(httpClientRequest.sink.add)
      .catchError(httpClientRequest.sink.addError)
      .whenComplete(httpClientRequest.sink.close));

  var response = await http.Client().send(httpClientRequest);

  response.headers.remove('transfer-encoding');

  if (response.headers['content-encoding'] == 'gzip') {
    response.headers.remove('content-encoding');
    response.headers.remove('content-length');
  }

  Content resp = new Content();

  resp.content = response.stream;
  resp.headers = response.headers;

  return resp;
}
