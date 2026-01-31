import 'package:flutter/material.dart';
import 'package:texas_holdem_frontend/api/client.dart';

class EvaluateTab extends StatefulWidget {
  const EvaluateTab({super.key});

  @override
  State<EvaluateTab> createState() => _EvaluateTabState();
}

class _EvaluateTabState extends State<EvaluateTab> {
  final _hole1 = TextEditingController(text: 'HA');
  final _hole2 = TextEditingController(text: 'HK');
  final _c1 = TextEditingController(text: 'HQ');
  final _c2 = TextEditingController(text: 'HJ');
  final _c3 = TextEditingController(text: 'HT');
  final _c4 = TextEditingController(text: 'S2');
  final _c5 = TextEditingController(text: 'D3');

  final _client = ApiClient();
  Map<String, dynamic>? _result;
  String? _error;
  bool _loading = false;

  Future<void> _evaluate() async {
    setState(() {
      _result = null;
      _error = null;
      _loading = true;
    });
    try {
      final r = await _client.evaluate(
        holeCards: [_hole1.text.trim(), _hole2.text.trim()],
        communityCards: [
          _c1.text.trim(),
          _c2.text.trim(),
          _c3.text.trim(),
          _c4.text.trim(),
          _c5.text.trim(),
        ],
      );
      setState(() {
        _result = r;
        _loading = false;
      });
    } on ApiException catch (e) {
      setState(() {
        _error = e.message;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  @override
  void dispose() {
    _hole1.dispose();
    _hole2.dispose();
    _c1.dispose();
    _c2.dispose();
    _c3.dispose();
    _c4.dispose();
    _c5.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          const Text(
            '2 hole cards + 5 community cards → best hand',
            style: TextStyle(fontSize: 16),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(child: _cardField('Hole 1', _hole1)),
              const SizedBox(width: 8),
              Expanded(child: _cardField('Hole 2', _hole2)),
            ],
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              Expanded(child: _cardField('Comm 1', _c1)),
              const SizedBox(width: 8),
              Expanded(child: _cardField('Comm 2', _c2)),
              const SizedBox(width: 8),
              Expanded(child: _cardField('Comm 3', _c3)),
              const SizedBox(width: 8),
              Expanded(child: _cardField('Comm 4', _c4)),
              const SizedBox(width: 8),
              Expanded(child: _cardField('Comm 5', _c5)),
            ],
          ),
          const SizedBox(height: 16),
          const Text('Format: 2 chars, e.g. HA (Heart Ace), S7 (Spade 7), CT (Club Ten)'),
          const SizedBox(height: 16),
          FilledButton(
            onPressed: _loading ? null : _evaluate,
            child: _loading
                ? const SizedBox(
                    height: 20,
                    width: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('Evaluate'),
          ),
          if (_error != null) ...[
            const SizedBox(height: 16),
            Card(
              color: Theme.of(context).colorScheme.errorContainer,
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Text(_error!),
              ),
            ),
          ],
          if (_result != null) ...[
            const SizedBox(height: 16),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Best hand: ${(_result!['best_hand'] as List).join(' ')}',
                      style: Theme.of(context).textTheme.titleMedium,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'Rank: ${_result!['rank_name']}',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                  ],
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _cardField(String label, TextEditingController c) {
    return TextField(
      controller: c,
      decoration: InputDecoration(
        labelText: label,
        border: const OutlineInputBorder(),
      ),
      maxLength: 2,
      textCapitalization: TextCapitalization.characters,
    );
  }
}
