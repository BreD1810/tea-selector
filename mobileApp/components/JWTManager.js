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
};

export default JWTManager;
