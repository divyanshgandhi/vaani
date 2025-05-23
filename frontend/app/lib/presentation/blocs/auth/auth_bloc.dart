import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:vaani/domain/repositories/auth_repository.dart';

// Events
abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class AuthStarted extends AuthEvent {
  const AuthStarted();
}

class PhoneNumberSubmitted extends AuthEvent {
  final String phoneNumber;

  const PhoneNumberSubmitted(this.phoneNumber);

  @override
  List<Object> get props => [phoneNumber];
}

class OtpSubmitted extends AuthEvent {
  final String otp;
  final String phoneNumber;

  const OtpSubmitted({required this.otp, required this.phoneNumber});

  @override
  List<Object> get props => [otp, phoneNumber];
}

class AuthLogoutRequested extends AuthEvent {
  const AuthLogoutRequested();
}

// States
abstract class AuthState extends Equatable {
  const AuthState();

  @override
  List<Object?> get props => [];
}

class AuthInitial extends AuthState {
  const AuthInitial();
}

class AuthLoading extends AuthState {
  const AuthLoading();
}

class PhoneNumberSubmitSuccess extends AuthState {
  const PhoneNumberSubmitSuccess();
}

class OtpVerificationInProgress extends AuthState {
  final String phoneNumber;

  const OtpVerificationInProgress(this.phoneNumber);

  @override
  List<Object> get props => [phoneNumber];
}

class AuthSuccess extends AuthState {
  const AuthSuccess();
}

class AuthFailure extends AuthState {
  final String message;

  const AuthFailure(this.message);

  @override
  List<Object> get props => [message];
}

// Bloc
class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _authRepository;

  AuthBloc({required AuthRepository authRepository})
      : _authRepository = authRepository,
        super(const AuthInitial()) {
    on<AuthStarted>(_onAuthStarted);
    on<PhoneNumberSubmitted>(_onPhoneNumberSubmitted);
    on<OtpSubmitted>(_onOtpSubmitted);
    on<AuthLogoutRequested>(_onAuthLogoutRequested);
  }

  Future<void> _onAuthStarted(
    AuthStarted event,
    Emitter<AuthState> emit,
  ) async {
    try {
      emit(const AuthLoading());
      final isLoggedIn = await _authRepository.isLoggedIn();

      if (isLoggedIn) {
        emit(const AuthSuccess());
      } else {
        emit(const AuthInitial());
      }
    } catch (e) {
      emit(AuthFailure(e.toString()));
    }
  }

  Future<void> _onPhoneNumberSubmitted(
    PhoneNumberSubmitted event,
    Emitter<AuthState> emit,
  ) async {
    try {
      emit(const AuthLoading());

      final success = await _authRepository.sendOtp(event.phoneNumber);

      if (success) {
        emit(OtpVerificationInProgress(event.phoneNumber));
      } else {
        emit(const AuthFailure('Failed to send OTP. Please try again.'));
      }
    } catch (e) {
      emit(AuthFailure(e.toString()));
    }
  }

  Future<void> _onOtpSubmitted(
    OtpSubmitted event,
    Emitter<AuthState> emit,
  ) async {
    try {
      emit(const AuthLoading());

      final success = await _authRepository.verifyOtp(
        event.phoneNumber,
        event.otp,
      );

      if (success) {
        emit(const AuthSuccess());
      } else {
        emit(const AuthFailure('Invalid OTP. Please try again.'));
      }
    } catch (e) {
      emit(AuthFailure(e.toString()));
    }
  }

  Future<void> _onAuthLogoutRequested(
    AuthLogoutRequested event,
    Emitter<AuthState> emit,
  ) async {
    try {
      await _authRepository.logout();
      emit(const AuthInitial());
    } catch (e) {
      emit(AuthFailure(e.toString()));
    }
  }
}
