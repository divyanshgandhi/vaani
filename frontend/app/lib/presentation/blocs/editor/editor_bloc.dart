import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../core/utils/text_formatter.dart';
import '../../../domain/models/script.dart';
import '../../../domain/repositories/pricing_repository.dart';
import '../../../domain/repositories/voice_repository.dart';
import 'editor_event.dart';
import 'editor_state.dart';

class EditorBloc extends Bloc<EditorEvent, EditorState> {
  final VoiceRepository _voiceRepository;
  final PricingRepository _pricingRepository;

  EditorBloc({
    required VoiceRepository voiceRepository,
    required PricingRepository pricingRepository,
  })  : _voiceRepository = voiceRepository,
        _pricingRepository = pricingRepository,
        super(const EditorState()) {
    on<TextChanged>(_onTextChanged);
    on<ApplyAutoPunctuation>(_onApplyAutoPunctuation);
    on<VoiceSelected>(_onVoiceSelected);
    on<EmotionChanged>(_onEmotionChanged);
    on<FetchVoices>(_onFetchVoices);
    on<FetchPricing>(_onFetchPricing);

    // Fetch initial data
    add(FetchVoices());
    add(FetchPricing());
  }

  void _onTextChanged(TextChanged event, Emitter<EditorState> emit) {
    final charCount = event.text.length;
    final script = state.script.copyWith(
      text: event.text,
      charCount: charCount,
    );
    final totalCost = charCount * state.pricePerChar;

    emit(state.copyWith(
      script: script,
      totalCost: totalCost,
      status: EditorStatus.ready,
    ));
  }

  void _onApplyAutoPunctuation(
    ApplyAutoPunctuation event,
    Emitter<EditorState> emit,
  ) {
    final formattedText = TextFormatter.autoPunctuate(state.script.text);
    final charCount = formattedText.length;
    final script = state.script.copyWith(
      text: formattedText,
      charCount: charCount,
    );
    final totalCost = charCount * state.pricePerChar;

    emit(state.copyWith(
      script: script,
      totalCost: totalCost,
      status: EditorStatus.ready,
    ));
  }

  void _onVoiceSelected(VoiceSelected event, Emitter<EditorState> emit) {
    final script = state.script.copyWith(
      voiceId: event.voiceId,
      voiceName: event.voiceName,
      language: event.language,
    );

    emit(state.copyWith(
      script: script,
      status: EditorStatus.ready,
    ));
  }

  void _onEmotionChanged(EmotionChanged event, Emitter<EditorState> emit) {
    final script = state.script.copyWith(
      emotionLevel: event.level,
    );

    emit(state.copyWith(
      script: script,
      status: EditorStatus.ready,
    ));
  }

  Future<void> _onFetchVoices(
    FetchVoices event,
    Emitter<EditorState> emit,
  ) async {
    try {
      emit(state.copyWith(status: EditorStatus.loadingVoices));
      final voices = await _voiceRepository.getVoices();

      // If no voice is selected, select the first voice
      Script script = state.script;
      if (state.script.voiceId == 'default' && voices.isNotEmpty) {
        script = script.copyWith(
          voiceId: voices.first.id,
          voiceName: voices.first.name,
          language: voices.first.language,
        );
      }

      emit(state.copyWith(
        availableVoices: voices,
        script: script,
        status: EditorStatus.ready,
      ));
    } catch (e) {
      emit(state.copyWith(
        status: EditorStatus.error,
        errorMessage: 'Failed to load voices: ${e.toString()}',
      ));
    }
  }

  Future<void> _onFetchPricing(
    FetchPricing event,
    Emitter<EditorState> emit,
  ) async {
    try {
      final pricePerChar = await _pricingRepository.getPricePerCharacter();
      final totalCost = state.script.charCount * pricePerChar;

      emit(state.copyWith(
        pricePerChar: pricePerChar,
        totalCost: totalCost,
      ));
    } catch (e) {
      // Just use the default price if we can't fetch from API
      emit(state.copyWith(
        pricePerChar: 0.01,
        totalCost: state.script.charCount * 0.01,
      ));
    }
  }
}
