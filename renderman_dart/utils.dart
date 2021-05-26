class Content {
  Map<String, dynamic> headers = {};
  dynamic content;

  Content(this.headers, this.content);

  factory Content.fromJson(dynamic json) {
    return Content(
        json['headers'] as Map<String, dynamic>, json['content'] as String);
  }

  Map toJson() {
    return {'content': content, 'headers': headers};
  }
}

bool isAssetRequest(String uri) {
  RegExp exp = new RegExp(
      r"(.*(css|js|json|jpg|jpeg|gif|svg|png|css|mp4|ogg|svg|ttf|ico|map))");

  return exp.hasMatch(uri);
}

bool isBotRequest(String agent) {
  RegExp exp = new RegExp(
      r"(googlebot|bingbot|adidxBot|yandex|baiduspider|twitterbot|facebookexternalhit|rogerbot|linkedinbot|embedly|quora link preview|showyoubot|outbrain|pinterest\/0\.|pinterestbot|slackbot|slack-imgproxy|vkshare|w3c_validator|whatsapp|collector-agent)");

  return exp.hasMatch(agent);
}
