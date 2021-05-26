import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'dart:async' show Future;
import 'package:shelf/shelf.dart' as shelf;
import 'package:shelf/shelf_io.dart' as io;
import 'package:dotenv/dotenv.dart' show load, env;
import 'package:sprintf/sprintf.dart';
import 'package:logging/logging.dart' show Logger, Level, LogRecord;
// import 'package:es_compression/zstd.dart';

import 'utils.dart';
import 'renderman.dart';
import 'neat_cache/neat_cache.dart';

final Logger log = Logger('dartis');
final cacheProvider = Cache.redisCacheProvider('redis://localhost:6379');
final cache = Cache(cacheProvider).withPrefix('prospeku-').withCodec(utf8);
// final codec = ZstdCodec();

void main() {
  load();

  Logger.root.level = Level.INFO;
  Logger.root.onRecord.listen((LogRecord rec) {
    print(
        '{"level": "${rec.level.name}", "time": "${rec.time}", "message": "${rec.message}" }');
  });

  var handler = const shelf.Pipeline()
      .addMiddleware(shelf.logRequests())
      .addHandler(Route);

  io.serve(handler, '0.0.0.0', int.parse(env['PORT'])).then((server) {
    print('Serving at http://${server.address.host}:${server.port}');
  });
}

Future<shelf.Response> Route(shelf.Request request) async {
  String uri = sprintf("%s/%s", [env['BACKEND_URL'], request.url]);
  String uriKey = sprintf("%s/%s", [env['APP_URL'], request.url]);
  String key = sprintf("%s", [md5.convert(utf8.encode(uriKey))]);
  Content resp;

  var agent = request.headers['user-agent'];

  try {
    if (isAssetRequest(uri)) {
      log.info({
        "agent": agent,
        "context": "Request as proxy request",
        "uri": uri,
      });

      resp = await fetchHttp(uri, request);
      return new shelf.Response.ok(
        resp.content,
        headers: resp.headers,
      );
    }

    if (isBotRequest(agent)) {
      log.info({
        "agent": agent,
        "context": "Request as bot request",
        "uri": uri,
      });

      var Rcontent = await cache[key].get();
      if (agent == 'collector-agent' || Rcontent == null) {
        resp = await fetchHeadless(uri, request);
        await cache[key].set(jsonEncode(resp), new Duration(days: 1));
      }

      Content content = Content.fromJson(jsonDecode(Rcontent));

      return new shelf.Response.ok(
        content.content,
        headers: content.headers,
      );
    }

    log.info({
      "agent": agent,
      "context": "Request as base request",
      "uri": uri,
    });

    resp = await fetchHttp(uri, request);
    return new shelf.Response.ok(
      resp.content,
      headers: resp.headers,
    );
  } catch (e) {
    log.severe(e);
    return new shelf.Response.internalServerError(body: "internal server");
  }
}
