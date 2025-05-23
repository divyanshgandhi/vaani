import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:vaani/domain/repositories/auth_repository.dart';
import 'package:vaani/presentation/blocs/auth/auth_bloc.dart';

class MockAuthRepository extends Mock implements AuthRepository {}

void main() {
  late AuthRepository authRepository;
  late AuthBloc authBloc;

  setUp(() {
    authRepository = MockAuthRepository();
    authBloc = AuthBloc(authRepository: authRepository);
  });

  tearDown(() {
    authBloc.close();
  });

  group('AuthBloc', () {
    test('initial state is AuthInitial', () {
      expect(authBloc.state, isA<AuthInitial>());
    });

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthInitial] when AuthStarted is added and user is not logged in',
      build: () {
        when(() => authRepository.isLoggedIn()).thenAnswer((_) async => false);
        return authBloc;
      },
      act: (bloc) => bloc.add(const AuthStarted()),
      expect: () => [
        const AuthLoading(),
        const AuthInitial(),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthSuccess] when AuthStarted is added and user is logged in',
      build: () {
        when(() => authRepository.isLoggedIn()).thenAnswer((_) async => true);
        return authBloc;
      },
      act: (bloc) => bloc.add(const AuthStarted()),
      expect: () => [
        const AuthLoading(),
        const AuthSuccess(),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, OtpVerificationInProgress] when PhoneNumberSubmitted is added and sendOtp succeeds',
      build: () {
        when(() => authRepository.sendOtp(any())).thenAnswer((_) async => true);
        return authBloc;
      },
      act: (bloc) => bloc.add(const PhoneNumberSubmitted('+911234567890')),
      expect: () => [
        const AuthLoading(),
        const OtpVerificationInProgress('+911234567890'),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthFailure] when PhoneNumberSubmitted is added and sendOtp fails',
      build: () {
        when(() => authRepository.sendOtp(any()))
            .thenAnswer((_) async => false);
        return authBloc;
      },
      act: (bloc) => bloc.add(const PhoneNumberSubmitted('+911234567890')),
      expect: () => [
        const AuthLoading(),
        const AuthFailure('Failed to send OTP. Please try again.'),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthSuccess] when OtpSubmitted is added and verifyOtp succeeds',
      build: () {
        when(() => authRepository.verifyOtp(any(), any()))
            .thenAnswer((_) async => true);
        return authBloc;
      },
      act: (bloc) => bloc
          .add(const OtpSubmitted(otp: '123456', phoneNumber: '+911234567890')),
      expect: () => [
        const AuthLoading(),
        const AuthSuccess(),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthLoading, AuthFailure] when OtpSubmitted is added and verifyOtp fails',
      build: () {
        when(() => authRepository.verifyOtp(any(), any()))
            .thenAnswer((_) async => false);
        return authBloc;
      },
      act: (bloc) => bloc
          .add(const OtpSubmitted(otp: '123456', phoneNumber: '+911234567890')),
      expect: () => [
        const AuthLoading(),
        const AuthFailure('Invalid OTP. Please try again.'),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthInitial] when AuthLogoutRequested is added',
      build: () {
        when(() => authRepository.logout()).thenAnswer((_) async {});
        return authBloc;
      },
      act: (bloc) => bloc.add(const AuthLogoutRequested()),
      expect: () => [
        const AuthInitial(),
      ],
    );

    blocTest<AuthBloc, AuthState>(
      'emits [AuthFailure] when an exception is thrown during logout',
      build: () {
        when(() => authRepository.logout())
            .thenThrow(Exception('Logout failed'));
        return authBloc;
      },
      act: (bloc) => bloc.add(const AuthLogoutRequested()),
      expect: () => [
        isA<AuthFailure>(),
      ],
    );
  });
}
