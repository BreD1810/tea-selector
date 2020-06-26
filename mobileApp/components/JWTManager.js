import {AsyncStorage} from 'react-native';

const JWTManager = {
  async getJWT(setJWT, setIsAuthorized) {
    try {
      const value = await AsyncStorage.getItem('jwt_token');
      if (value !== null) {
        console.log(`JWT: ${value}`);
        console.log(value);
        setJWT(value);
        setIsAuthorized(true);
      } else {
        console.log('No JWT');
      }
    } catch (error) {
      console.warn('AsyncStorage error: ' + error.message);
    }
  },
  async setJWT(jwt) {
    try {
      await AsyncStorage.setItem('jwt_token', jwt);
    } catch (error) {
      console.warn('AsyncStorage error: ' + error.message);
    }
  },
  async clearJWT() {
    try{
      await AsyncStorage.removeItem('jwt_token');
    } catch (error) {
      console.warn('AsynStorage error: ' + error.message);
    }
  },
};

export default JWTManager;
