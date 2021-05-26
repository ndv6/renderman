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
      r"(.*(js|css|xml|less|png|jpg|jpeg|gif|pdf|doc|txt|ico|rss|zip|mp3|rar|exe|wmv|doc|avi|ppt|mpg|mpeg|tif|wav|mov|psd|ai|xls|mp4|m4a|swf|dat|dmg|iso|flv|m4v|torrent|ttf|woff|svg|eot))");

  return exp.hasMatch(uri);
}

bool isBotRequest(String agent) {
  RegExp exp = new RegExp(
      r"(googlebot|bingbot|adidxBot|yandex|baiduspider|twitterbot|facebookexternalhit|rogerbot|linkedinbot|embedly|quora link preview|showyoubot|outbrain|pinterest\/0\.|pinterestbot|slackbot|slack-imgproxy|vkshare|w3c_validator|whatsapp|collector-agent)");

  return exp.hasMatch(agent);
}
