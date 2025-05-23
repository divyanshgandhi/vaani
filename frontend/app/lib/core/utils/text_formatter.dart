/// A utility class for text formatting operations
class TextFormatter {
  /// Applies auto-punctuation to improve readability of the given text
  ///
  /// This is a simple implementation that adds or fixes common punctuation issues.
  /// In a real app, this would potentially use a more sophisticated NLP model.
  static String autoPunctuate(String text) {
    if (text.isEmpty) return text;

    // Capitalize first letter of sentences
    String processed = text.replaceAllMapped(RegExp(r'(^|[.!?]\s+)([a-z])'),
        (match) => '${match[1]}${match[2]!.toUpperCase()}');

    // Add period at the end if missing
    if (!processed.endsWith('.') &&
        !processed.endsWith('!') &&
        !processed.endsWith('?')) {
      processed = '$processed.';
    }

    // Fix spaces before punctuation
    processed = processed.replaceAll(RegExp(r'\s+([.,!?:;])'), r'$1');

    // Ensure space after punctuation
    processed = processed.replaceAllMapped(RegExp(r'([.,!?:;])([a-zA-Z0-9])'),
        (match) => '${match[1]} ${match[2]}');

    // Remove double spaces
    processed = processed.replaceAll(RegExp(r'\s{2,}'), ' ');

    return processed;
  }
}
