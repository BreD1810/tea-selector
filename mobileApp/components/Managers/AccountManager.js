import React, {useState} from 'react';
import {
  Alert,
  TextInput,
  ToastAndroid,
  Text,
  KeyboardAvoidingView,
  ScrollView,
  StyleSheet,
  Dimensions,
} from 'react-native';
import Button from 'react-native-button';
import {serverURL} from '../../app.json';
import JWTManager from '../JWTManager';

const AccountManager = ({jwtToken, setLoggedIn}) => {
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [newPasswordRepeat, setNewPasswordRepeat] = useState('');

  const attemptChangePassword = () => {
    if (oldPassword === '') {
      Alert.alert('Error changing password', 'Please enter your old password');
      return;
    }
    if (newPassword === '') {
      Alert.alert('Error changing password', 'Please enter your new password');
      return;
    }
    if (newPasswordRepeat === '') {
      Alert.alert('Error changing password', 'Please repeat your new password');
      return;
    }
    if (newPassword !== newPasswordRepeat) {
      Alert.alert('Error changing password', 'New passwords do not match');
      return;
    }
    if (oldPassword === newPassword) {
      Alert.alert(
        'Error changing password',
        'Your new password cannot be your old password',
      );
      return;
    }

    fetch(serverURL + '/changepassword', {
      method: 'POST',
      body: JSON.stringify({old: oldPassword, new: newPassword}),
      headers: {
        Token: jwtToken,
      },
    })
      .then(response => response.json())
      .then(json => {
        if (json.error) {
          throw new Error(json.error);
        }
        return json;
      })
      .then(json => {
        ToastAndroid.showWithGravityAndOffset(
          'Successfully changed password!',
          ToastAndroid.SHORT,
          ToastAndroid.BOTTOM,
          0,
          150,
        );
        this.oldPasswordInput.clear();
        this.newPasswordInput.clear();
        this.newPasswordRepeatInput.clear();
        return json;
      })
      .catch(error => {
        if (error.message === 'Incorrect password') {
          Alert.alert('Error logging in', 'Your old password is incorrect');
        } else {
          console.warn(error);
          Alert.alert('Error logging in', 'Please try again.');
        }
      });
  };

  const logout = () => {
    JWTManager.clearJWT();
    ToastAndroid.showWithGravityAndOffset(
      'Successfully logged out!',
      ToastAndroid.SHORT,
      ToastAndroid.BOTTOM,
      0,
      150,
    );
    setLoggedIn(false);
  };

  return (
    <ScrollView>
      <KeyboardAvoidingView behavior={'height'} style={styles.container}>
        <Text style={styles.title}>Change Password</Text>
        <TextInput
          placeholder={'Old Password'}
          secureTextEntry={true}
          style={styles.input}
          onChangeText={textValue => setOldPassword(textValue)}
          selection="center"
          ref={input => {
            this.oldPasswordInput = input;
          }}
        />
        <TextInput
          placeholder={'New Password'}
          secureTextEntry={true}
          style={styles.input}
          onChangeText={textValue => setNewPassword(textValue)}
          selection="center"
          ref={input => {
            this.newPasswordInput = input;
          }}
        />
        <TextInput
          placeholder={'Repeat New Password'}
          secureTextEntry={true}
          style={styles.input}
          onChangeText={textValue => setNewPasswordRepeat(textValue)}
          selection="center"
          ref={input => {
            this.newPasswordRepeatInput = input;
          }}
        />
        <Button
          onPress={() => attemptChangePassword()}
          containerStyle={styles.btnContainer}
          style={styles.btn}>
          Change Password
        </Button>
        <Text style={styles.title2}>Log Out</Text>
        <Button
          onPress={() => logout()}
          containerStyle={styles.btnContainerBottom}
          style={styles.btn}>
          Log Out
        </Button>
      </KeyboardAvoidingView>
    </ScrollView>
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
  title: {
    fontSize: 32,
    fontWeight: '600',
    color: 'black',
  },
  title2: {
    fontSize: 32,
    fontWeight: '600',
    color: 'black',
    marginTop: 40,
  },
  input: {
    height: 40,
    padding: 5,
    marginTop: 20,
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
    marginTop: 25,
    height: 60,
    width: 200,
    overflow: 'hidden',
    borderRadius: 4,
    backgroundColor: 'dodgerblue',
  },
  btnContainerBottom: {
    padding: 15,
    marginTop: 25,
    marginBottom: 80,
    height: 60,
    width: 200,
    overflow: 'hidden',
    borderRadius: 4,
    backgroundColor: 'dodgerblue',
  },
});

export default AccountManager;
