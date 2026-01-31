import 'package:flutter/material.dart';
import 'package:texas_holdem_frontend/api/client.dart';

class ProbabilityTab extends StatefulWidget {
  const ProbabilityTab({super.key});

  @override
  State<ProbabilityTab> createState() => _ProbabilityTabState();
}

class _ProbabilityTabState extends State<ProbabilityTab> {
  final _h1 = TextEditingController(text: 'HA');
  final _h2 = TextEditingController(text: 'HK');
  final _c1 = TextEditingController();
  final _c2 = TextEditingController();
  final _c3 = TextEditingController();
  final _c4 = TextEditingController();
  final _c5 = TextEditingController();
  final _numPlayersController = TextEditingController(text: '2');
  final _numSimsController = TextEditingController(text: '10000');

  final _client = ApiClient();
  Map<String, dynamic>? _result;
  String? _error;
  bool _loading = false;

  List<String> get _communityCards {
    final list = <String>[];
    for (var c in [_c1, _c2, _c3, _c4, _c5]) {
      final t = c.text.trim();
      if (t.isNotEmpty) list.add(t);
    }
    return list;
  }

  Future<void> _calculate() async {
    setState(() {
      _result = null;
      _error = null;
      _loading = true;
    });
    try {
      final r = await _client.probability(
        holeCards: [_h1.text.trim(), _h2.text.trim()],
        communityCards: _communityCards,
        numPlayers: int.tryParse(_numPlayersController.text) ?? 2,
        numSims: int.tryParse(_numSimsController.text) ?? 10000,
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
    _numPlayersController.dispose();
    _numSimsController.dispose();
    _h1.dispose();
    _h2.dispose();
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
            'Win probability via Monte Carlo (0–5 community cards)',
            style: TextStyle(fontSize: 16),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(child: _cardField('Hole 1', _h1)),
              const SizedBox(width: 8),
              Expanded(child: _cardField('Hole 2', _h2)),
            ],
          ),
          const SizedBox(height: 8),
          const Text('Community cards (0–5, leave empty for pre-flop)'),
          Row(
            children: [
              for (var c in [_c1, _c2, _c3, _c4, _c5])
                Expanded(child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 4),
                  child: TextField(controller: c, maxLength: 2, textCapitalization: TextCapitalization.characters),
                )),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              const Text('Number of players: '),
              SizedBox(
                width: 80,
                child: TextField(
                  controller: _numPlayersController,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(border: OutlineInputBorder()),
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              const Text('Monte Carlo simulations: '),
              SizedBox(
                width: 100,
                child: TextField(
                  controller: _numSimsController,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(border: OutlineInputBorder()),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          FilledButton(
            onPressed: _loading ? null : _calculate,
            child: _loading
                ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2))
                : const Text('Calculate Probability'),
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
                    Text(
                      'Win probability: ${((_result!['win_probability'] as num) * 100).toStringAsFixed(2)}%',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    Text('(based on ${_result!['num_sims']} simulations, ${_result!['num_players']} players)'),
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
      decoration: InputDecoration(labelText: label, border: const OutlineInputBorder()),
      maxLength: 2,
      textCapitalization: TextCapitalization.characters,
    );
  }
}
