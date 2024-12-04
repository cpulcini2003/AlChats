<template>
  <div id="app">
    <nav class="navbar">
      <!-- <div class="navbar-title">Chats</div> -->
      
      <!-- Left-aligned buttons -->
      <div class="navbar-left">
        <router-link to="/" class="navbar-button">
          <button>Home</button>
        </router-link>
        <router-link to="/users" class="navbar-button">
          <button>Users</button>
        </router-link>
        <router-link to="/about" class="navbar-button">
          <button>About</button>
        </router-link>
      </div>


      <!-- Right-aligned login/logout button -->
      <div class="navbar-right">
        <div class="navbar-title" v-if="username">
          {{ username }}
        </div>
        <router-link v-if="!username" to="/login" class="navbar-button">
          <button>Login</button>
        </router-link>
        <router-link v-if="username" to="/" class="navbar-button">
          <button @click="logout">Logout</button>
        </router-link>

      </div>
    </nav>

    <!-- Routed components will be rendered here -->
    <router-view></router-view>
  </div>
</template>

<script>
// import { mapState } from 'vuex';

export default {
  name: "App",
  // computed: {
  //   ...mapState(['username']), // Map the username from Vuex store to this component
  // },
  methods: {
    logout() {
      // Clear username in Vuex store and sessionStorage
      this.$store.commit("setUsername", null); // Reset username in Vuex store
    },
  },
};
</script>

<style>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #1ce036;
  padding: 1em;
}

.navbar-title {
  color: white;
  font-size: 1.5em;
  font-weight: bold;
}


.navbar-buttons {
  display: flex;
  gap: 1em;
}

.navbar-button button {
  background-color: white;
  color: #2c3e50;
  padding: 10px 15px;
  border: 1px solid white;
  border-radius: 4px;
  font-size: 1em;
  cursor: pointer;
  transition: background-color 0.3s, color 0.3s;
}

.navbar-button button:hover {
  background-color: #34495e;
  color: white;
}

/* Add padding or margin between buttons */
.navbar-left .navbar-button {
  margin-right: 10px; /* Space between buttons */
}

.navbar-left .navbar-button:last-child {
  margin-right: 0; /* Remove margin for the last button */
}

.navbar-button button {
  padding: 8px 16px; /* Optional: Style individual buttons */
}

/* Style the right section as a flexbox */
.navbar-right {
  display: flex;
  align-items: center; /* Align items vertically */
  gap: 10px; /* Add spacing between UNO and DUE */
}


</style>
