import 'package:flutter/material.dart';
import 'package:texas_holdem_frontend/screens/home_screen.dart';

void main() {
  runApp(const TexasHoldemApp());
}

class TexasHoldemApp extends StatelessWidget {
  const TexasHoldemApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Texas Hold\'em',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: Colors.green.shade800,
          brightness: Brightness.dark,
        ),
        useMaterial3: true,
      ),
      home: const HomeScreen(),
    );
  }
}
