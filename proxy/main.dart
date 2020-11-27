import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:shelf/shelf.dart' as shelf;
import 'package:shelf/shelf_io.dart' as io;
import 'package:dotenv/dotenv.dart' show load, env;
import 'package:sprintf/sprintf.dart';
import 'package:dcache/dcache.dart';

import 'utils.dart';
import 'headless.dart';
import 'http.dart';

Cache cache;

void main() async {
  load();

  cache = new SimpleCache(storage: new SimpleStorage(size: 20));

  var handler = const shelf.Pipeline()
      .addMiddleware(shelf.logRequests())
      .addHandler(echoRequest);

  io.serve(handler, '0.0.0.0', int.parse(env['PORT'])).then((server) {
    print('Serving at http://${server.address.host}:${server.port}');
  });
}

Future<shelf.Response> echoRequest(shelf.Request request) async {
  String uri = sprintf("%s/%s", [env['BACKEND_URL'], request.url]);
  String hash = sprintf("%s", [md5.convert(utf8.encode(uri))]);
  Content resp = Content();

  if (isAssetRequest(uri)) {
    resp = await fetchHttp(uri, request);
  } else if (isBotRequest(request.headers['user-agent'])) {
    Content c = cache.get(hash);
    if (c != null) {
      resp = c;
    } else {
      resp = await fetchHeadless(uri, request);
      cache.set(hash, resp);
    }
  } else {
    resp = await fetchHttp(uri, request);
  }

  return new shelf.Response.ok(resp.content, headers: resp.headers);
}
