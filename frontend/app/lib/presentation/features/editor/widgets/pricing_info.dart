import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';
import '../../../blocs/editor/editor_bloc.dart';
import '../../../blocs/editor/editor_state.dart';

class PricingInfo extends StatelessWidget {
  const PricingInfo({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final theme = Theme.of(context);

    return BlocBuilder<EditorBloc, EditorState>(
      builder: (context, state) {
        return Card(
          elevation: 0,
          color: theme.colorScheme.surfaceVariant,
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Row(
                      children: [
                        Icon(
                          Icons.info_outline,
                          size: 16,
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          l10n.estimatedCost(state.totalCost),
                          style: theme.textTheme.bodyMedium,
                        ),
                      ],
                    ),
                    Chip(
                      label: Text(
                        state.isOverCharLimit
                            ? 'Limit Exceeded'
                            : '${state.script.charCount}/1000',
                        style: TextStyle(
                          color: state.isOverCharLimit
                              ? theme.colorScheme.onError
                              : theme.colorScheme.onSurface,
                        ),
                      ),
                      backgroundColor: state.isOverCharLimit
                          ? theme.colorScheme.error
                          : null,
                      padding: const EdgeInsets.all(0),
                    ),
                  ],
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
