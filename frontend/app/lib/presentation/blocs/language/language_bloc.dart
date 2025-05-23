import 'dart:ui';

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:hydrated_bloc/hydrated_bloc.dart';

// Events
abstract class LanguageEvent extends Equatable {
  const LanguageEvent();

  @override
  List<Object?> get props => [];
}

class LanguageChanged extends LanguageEvent {
  final Locale locale;

  const LanguageChanged(this.locale);

  @override
  List<Object> get props => [locale];
}

// State
class LanguageState extends Equatable {
  final Locale locale;

  const LanguageState({required this.locale});

  factory LanguageState.initial() =>
      const LanguageState(locale: Locale('en', ''));

  LanguageState copyWith({Locale? locale}) {
    return LanguageState(locale: locale ?? this.locale);
  }

  @override
  List<Object> get props => [locale];

  Map<String, dynamic> toJson() {
    return {
      'languageCode': locale.languageCode,
      'countryCode': locale.countryCode ?? '',
    };
  }

  factory LanguageState.fromJson(Map<String, dynamic> json) {
    return LanguageState(
      locale: Locale(
        json['languageCode'] as String,
        json['countryCode'] as String,
      ),
    );
  }
}

// Bloc
class LanguageBloc extends HydratedBloc<LanguageEvent, LanguageState> {
  LanguageBloc() : super(LanguageState.initial()) {
    on<LanguageChanged>(_onLanguageChanged);
  }

  void _onLanguageChanged(LanguageChanged event, Emitter<LanguageState> emit) {
    emit(state.copyWith(locale: event.locale));
  }

  @override
  LanguageState? fromJson(Map<String, dynamic> json) {
    return LanguageState.fromJson(json);
  }

  @override
  Map<String, dynamic>? toJson(LanguageState state) {
    return state.toJson();
  }
}
