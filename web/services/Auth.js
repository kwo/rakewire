/* eslint prefer-template: 0 */
class Auth {

	static credentialsKey = 'credentials';

	static login(username, password) {
		return new Promise((resolve) => {
			fetch('/api/token', {
				method: 'POST',
				headers: {
					Authorization: 'Basic ' + btoa(`${username}:${password}`),
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({})
			}).then(rsp => {
				if (rsp.status === 200) {
					Auth.credentials = rsp.json();
					resolve(true);
				} else {
					resolve(false);
				}
			}).catch((/*x*/) => {
				//console.log('login failed:', x);
				resolve(false);
			});
		});
	}

	static logout() {
		Auth.credentials = null;
		return Promise.resolve(true);
	}

	static get loggedIn() {
		return !!Auth.credentials;
	}

	static get credentials() {
		const value = localStorage.getItem(Auth.credentialsKey);
		if (value) {
			return JSON.parse(value);
		}
		return null;
	}

	static set credentials(value) {
		if (value) {
			localStorage.setItem(Auth.credentialsKey, JSON.stringify(value));
		} else {
			localStorage.removeItem(Auth.credentialsKey);
		}
	}

}

export default Auth;
