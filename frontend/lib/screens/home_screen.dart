import 'package:flutter/material.dart';
import 'package:texas_holdem_frontend/screens/evaluate_tab.dart';
import 'package:texas_holdem_frontend/screens/compare_tab.dart';
import 'package:texas_holdem_frontend/screens/probability_tab.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Texas Hold\'em - Hand Evaluator'),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: 'Evaluate Hand'),
            Tab(text: 'Compare Hands'),
            Tab(text: 'Win Probability'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: const [
          EvaluateTab(),
          CompareTab(),
          ProbabilityTab(),
        ],
      ),
    );
  }
}
