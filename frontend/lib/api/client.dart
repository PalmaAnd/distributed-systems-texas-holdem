import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;

class ApiClient {
  ApiClient({
    String? baseUrl,
  }) : baseUrl = baseUrl ?? (kReleaseMode ? '' : 'http://localhost:8080');
  final String baseUrl;

  Future<Map<String, dynamic>> evaluate({
    required List<String> holeCards,
    required List<String> communityCards,
  }) async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/v1/evaluate'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'hole_cards': holeCards,
        'community_cards': communityCards,
      }),
    );
    return _handleResponse(res);
  }

  Future<Map<String, dynamic>> compare({
    required List<String> hole1,
    required List<String> community1,
    required List<String> hole2,
    required List<String> community2,
  }) async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/v1/compare'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'hand1': {'hole_cards': hole1, 'community_cards': community1},
        'hand2': {'hole_cards': hole2, 'community_cards': community2},
      }),
    );
    return _handleResponse(res);
  }

  Future<Map<String, dynamic>> probability({
    required List<String> holeCards,
    List<String> communityCards = const [],
    int numPlayers = 2,
    int numSims = 10000,
  }) async {
    final res = await http.post(
      Uri.parse('$baseUrl/api/v1/probability'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'hole_cards': holeCards,
        'community_cards': communityCards,
        'num_players': numPlayers,
        'num_sims': numSims,
      }),
    );
    return _handleResponse(res);
  }

  Map<String, dynamic> _handleResponse(http.Response res) {
    final body = jsonDecode(res.body) as Map<String, dynamic>;
    if (res.statusCode >= 400) {
      throw ApiException(
        body['error'] as String? ?? 'Request failed',
        res.statusCode,
      );
    }
    return body;
  }
}

class ApiException implements Exception {
  ApiException(this.message, this.statusCode);
  final String message;
  final int statusCode;
  @override
  String toString() => 'ApiException: $message ($statusCode)';
}
