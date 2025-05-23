import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../../../blocs/editor/editor_bloc.dart';
import '../../../blocs/editor/editor_event.dart';
import '../../../blocs/editor/editor_state.dart';

class EmotionSlider extends StatelessWidget {
  const EmotionSlider({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);

    return BlocBuilder<EditorBloc, EditorState>(
      buildWhen: (previous, current) =>
          previous.script.emotionLevel != current.script.emotionLevel,
      builder: (context, state) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  l10n.emotionIntensity,
                  style: theme.textTheme.titleMedium,
                ),
                Text(
                  '${state.script.emotionLevel}%',
                  style: theme.textTheme.bodySmall,
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                const Icon(Icons.sentiment_neutral, size: 20),
                Expanded(
                  child: Slider(
                    value: state.script.emotionLevel.toDouble(),
                    min: 0,
                    max: 100,
                    divisions: 10,
                    onChanged: (value) {
                      context
                          .read<EditorBloc>()
                          .add(EmotionChanged(value.toInt()));
                    },
                  ),
                ),
                const Icon(Icons.sentiment_very_satisfied, size: 20),
              ],
            ),
          ],
        );
      },
    );
  }
}
