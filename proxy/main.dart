import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:shelf/shelf.dart' as shelf;
import 'package:shelf/shelf_io.dart' as io;
import 'package:puppeteer/puppeteer.dart';
import 'package:dotenv/dotenv.dart' show load, env;
import 'package:sprintf/sprintf.dart';
import 'package:http/http.dart' as http;
import 'package:dcache/dcache.dart';

Browser browser;
Page page;
Cache cache;

class Content {
  Map<String, String> headers;
  String content;
}

var skippedResources = [
  'googletagmanager',
  'google-analytics',
  'facebook',
  'twitter',
  'windows',
  'service-worker'
];

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

  Content c = cache.get(hash);
  if (c != null) {
    resp.headers = c.headers;
    resp.content = c.content;
    return new shelf.Response.ok(resp.content, headers: resp.headers);
  }

  if (isBotRequest(request.headers['user-agent'])) {
    if (isAssetRequest(uri)) {
      http.Response response = await http.get(uri);
      resp.headers = response.headers;
      resp.content = response.body;
    } else {
      await setup();

      await page.goto(uri, wait: Until.networkIdle);
      resp.headers = {'content-type': 'text/html'};
      resp.content = await page.content;

      await browser.disconnect();
    }
  } else {
    http.Response response = await http.get(uri);
    resp.headers = response.headers;
    resp.content = response.body;
  }

  if (c == null) {
    cache.set(hash, resp);
  }

  return new shelf.Response.ok(resp.content, headers: resp.headers);
}

void setup() async {
  if (env['CHROMIUM_URL'] != null) {
    browser = await puppeteer.connect(browserWsEndpoint: env['CHROMIUM_URL']);
  } else {
    browser = await puppeteer.launch(
      // headless: false,
      // devTools: true,
      ignoreHttpsErrors: true,
      noSandboxFlag: true,
      userDataDir: './tmp/chrome',
      args: [
        '--disable-web-security',
        '--no-sandbox',
        '--disable-setuid-sandbox',
        '--disable-dev-shm-usage',
        '--disable-accelerated-2d-canvas',
        '--disable-gpu',
        '--disable-notifications',
      ],
    );
  }

  page = await browser.newPage();
  await page.setRequestInterception(true);
  page.onRequest.listen((Request request) {
    for (String item in skippedResources) {
      if (request.url.contains(item)) {
        request.abort(error: ErrorReason.blockedByClient);
        break;
      }
    }

    if (request.resourceType == ResourceType.image) {
      request.abort(error: ErrorReason.blockedByClient);
    } else {
      request.continueRequest(headers: request.headers);
    }
  });
}

bool isAssetRequest(String uri) {
  RegExp exp = new RegExp(
      r"(.*(css|js|json|jpg|jpeg|gif|svg|png|css|mp4|ogg|svg|ttf|ico|map))");

  return exp.hasMatch(uri);
}

bool isBotRequest(String agent) {
  RegExp exp = new RegExp(
      r"(googlebot|bingbot|yandex|twitterbot|facebookexternalhit|linkedinbot|embedly|pinterest\/0\.|pinterestbot|slackbot|whatsapp|curl)");

  return exp.hasMatch(agent);
}
