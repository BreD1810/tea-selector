import React, {useState} from 'react';
import {
  TextInput,
  KeyboardAvoidingView,
  Text,
  StyleSheet,
  Dimensions,
} from 'react-native';
import Button from 'react-native-button';

const Login = ({setJWT}) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const attemptLogin = () => {
    console.log('Attempting login');
  };

  return (
    <KeyboardAvoidingView behavior={'height'} style={styles.container}>
      <Text style={styles.title}>Login</Text>
      <TextInput
        placeholder={'Username'}
        style={styles.input}
        onChangeText={textValue => setUsername(textValue)}
        selection="center"
        ref={input => {
          this.usernameInput = input;
        }}
      />
      <TextInput
        placeholder={'Password'}
        style={styles.input}
        onChangeText={textValue => setPassword(textValue)}
        selection="center"
        ref={input => {
          this.passwordInput = input;
        }}
      />
      <Button
        onPress={attemptLogin}
        containerStyle={styles.btnContainer}
        style={styles.btn}>
        Login
      </Button>
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  container: {
    paddingTop: 30,
    flex: 1,
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
  },
  teaTextContainer: {
    height: 40,
  },
  title: {
    fontSize: 32,
    fontWeight: '600',
    color: 'black',
  },
  input: {
    height: 40,
    padding: 5,
    marginTop: 50,
    fontSize: 18,
    borderRadius: 10,
    width: Math.round(Dimensions.get('window').width) - 50,
    borderColor: 'lightgray',
    borderWidth: 2,
    textAlign: 'center',
  },
  btn: {
    fontSize: 20,
    color: 'white',
  },
  btnContainer: {
    padding: 15,
    marginTop: 50,
    marginBottom: 40,
    height: 60,
    width: 200,
    overflow: 'hidden',
    borderRadius: 4,
    backgroundColor: 'dodgerblue',
  },
});

export default Login;
