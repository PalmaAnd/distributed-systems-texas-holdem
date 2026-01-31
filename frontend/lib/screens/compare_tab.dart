import 'package:flutter/material.dart';
import 'package:texas_holdem_frontend/api/client.dart';

class CompareTab extends StatefulWidget {
  const CompareTab({super.key});

  @override
  State<CompareTab> createState() => _CompareTabState();
}

class _CompareTabState extends State<CompareTab> {
  final _h1a = TextEditingController(text: 'HA');
  final _h1b = TextEditingController(text: 'HK');
  final _c1a = TextEditingController(text: 'HQ');
  final _c1b = TextEditingController(text: 'HJ');
  final _c1c = TextEditingController(text: 'HT');
  final _c1d = TextEditingController(text: 'S2');
  final _c1e = TextEditingController(text: 'D3');

  final _h2a = TextEditingController(text: 'C2');
  final _h2b = TextEditingController(text: 'C3');
  final _c2a = TextEditingController(text: 'C4');
  final _c2b = TextEditingController(text: 'C5');
  final _c2c = TextEditingController(text: 'C6');
  final _c2d = TextEditingController(text: 'S7');
  final _c2e = TextEditingController(text: 'D8');

  final _client = ApiClient();
  Map<String, dynamic>? _result;
  String? _error;
  bool _loading = false;

  Future<void> _compare() async {
    setState(() {
      _result = null;
      _error = null;
      _loading = true;
    });
    try {
      final r = await _client.compare(
        hole1: [_h1a.text.trim(), _h1b.text.trim()],
        community1: [_c1a.text.trim(), _c1b.text.trim(), _c1c.text.trim(), _c1d.text.trim(), _c1e.text.trim()],
        hole2: [_h2a.text.trim(), _h2b.text.trim()],
        community2: [_c2a.text.trim(), _c2b.text.trim(), _c2c.text.trim(), _c2d.text.trim(), _c2e.text.trim()],
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
    _h1a.dispose();
    _h1b.dispose();
    _c1a.dispose();
    _c1b.dispose();
    _c1c.dispose();
    _c1d.dispose();
    _c1e.dispose();
    _h2a.dispose();
    _h2b.dispose();
    _c2a.dispose();
    _c2b.dispose();
    _c2c.dispose();
    _c2d.dispose();
    _c2e.dispose();
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
            'Compare two hands (2 hole + 5 community each)',
            style: TextStyle(fontSize: 16),
          ),
          const SizedBox(height: 16),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Hand 1', style: TextStyle(fontWeight: FontWeight.bold)),
                  Row(children: [
                    Expanded(child: _cardField('H1a', _h1a)),
                    const SizedBox(width: 4),
                    Expanded(child: _cardField('H1b', _h1b)),
                  ]),
                  Row(children: [
                    for (var c in [_c1a, _c1b, _c1c, _c1d, _c1e])
                      Expanded(child: _cardField('', c)),
                  ]),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Hand 2', style: TextStyle(fontWeight: FontWeight.bold)),
                  Row(children: [
                    Expanded(child: _cardField('H2a', _h2a)),
                    const SizedBox(width: 4),
                    Expanded(child: _cardField('H2b', _h2b)),
                  ]),
                  Row(children: [
                    for (var c in [_c2a, _c2b, _c2c, _c2d, _c2e])
                      Expanded(child: _cardField('', c)),
                  ]),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          FilledButton(
            onPressed: _loading ? null : _compare,
            child: _loading
                ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2))
                : const Text('Compare'),
          ),
          if (_error != null) ...[
            const SizedBox(height: 16),
            Card(
              color: Theme.of(context).colorScheme.errorContainer,
              child: Padding(padding: const EdgeInsets.all(16), child: Text(_error!)),
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
                    Text('Winner: ${_result!['winner']}', style: Theme.of(context).textTheme.titleLarge),
                    const SizedBox(height: 8),
                    Text('Hand 1: ${(_result!['hand1']['best_hand'] as List).join(' ')} → ${_result!['hand1']['rank_name']}'),
                    Text('Hand 2: ${(_result!['hand2']['best_hand'] as List).join(' ')} → ${_result!['hand2']['rank_name']}'),
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
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: TextField(
        controller: c,
        decoration: InputDecoration(labelText: label.isEmpty ? null : label, border: const OutlineInputBorder()),
        maxLength: 2,
        textCapitalization: TextCapitalization.characters,
      ),
    );
  }
}
