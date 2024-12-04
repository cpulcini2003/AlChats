<template>
  <div class="login-page">
    <h2>Login</h2>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">Username:</label>
        <input
          type="text"
          id="username"
          v-model="username"
          placeholder="Enter your username"
          required
        />
      </div>

      <div v-if="error" class="error-message">{{ error }}</div>

      <button type="submit" :disabled="isLoading">
        {{ isLoading ? "Logging in..." : "Login" }}
      </button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';
// import { mapActions } from 'vuex'; // Import mapActions to map store actions

export default {
  name: "LoginPage",
  data() {
    return {
      username: "",
      isLoading: false,
      error: null,
    };
  },
  methods: {
    // ...mapActions(['login']), // Map login action to this component

    async handleLogin() {
      this.error = null; // Reset error
      if (!this.username) {
        this.error = "Please enter a username.";
        return;
      }

      this.isLoading = true;

      try {
        const response = await axios.post("http://localhost:8082/user/session", null, {
          params: { username: this.username },
        });

        if (response.status === 200 || response.status === 201) {
          // // Successful login, save the username in sessionStorage
          // sessionStorage.setItem("username", this.username);
          
          // Update Vuex store with username
          this.login(this.username); // Dispatch login action to Vuex store

          // Redirect to users page after login
          this.$router.push("/users");
        } else {
          this.error = "Login failed. Please try again.";
        }
      } catch (err) {
        this.error = err.response?.data?.message || "Failed to log in. Please try again.";
      } finally {
        this.isLoading = false;
      }
    },
  },

};
</script>


<style scoped>
.login-page {
  max-width: 400px;
  margin: 50px auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #f9f9f9;
}

h2 {
  text-align: center;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

input {
  width: 100%;
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background-color: #c0c0c0;
}

.error-message {
  color: red;
  margin-bottom: 10px;
  font-size: 14px;
}
</style>
