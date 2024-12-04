<template>
  <div>
    <div v-if="username" class="user-list">
      <div
        class="user-card"
        v-for="user in users"
        :key="user.userId"
      >
        <!-- Display user photo, name, and userId -->
        <!-- <img :src="user.photo" alt="User Photo" class="user-photo" /> -->
        <h3>{{ user.username }}</h3>
        <p>User ID: {{ user.userId }}</p>
      </div>
    </div>
    <div v-if="!username" class="center-screen">Please login</div>
  </div>
</template>

<script>
import axios from "axios";
// import { mapState } from 'vuex';

export default {
  name: "UserList",
  // computed: {
  //   ...mapState(['username']), // Map the username from Vuex store to this component
  // },
  data() {
    return {
      users: [], // To store the user data
      loading: true, // To track loading state
      error: null, // To track errors
    };
  },
  mounted() {
    this.fetchUsers(); // Fetch data when the component is mounted
    this.interval = setInterval(this.fetchUsers, 1000);
  },
  methods: {
    async fetchUsers() {
      try {
        const response = await axios.get("http://localhost:8082/user"); // Replace with your API endpoint
        // Assuming the structure of the data returned is the same
        // Adjust the mapping here if the data needs to be transformed
        this.users = response.data.map(user => ({
          userId: user.userId, // Adjust according to your API response
          username: user.username, // Adjust according to your API response
          photo: user.photo || "https://example.com/default-photo.jpg" // Default photo if not available
        }));
        this.loading = false; // Set loading to false once the data is fetched
      } catch (err) {
        this.error = "Failed to load user data."; // Handle error if any
        this.loading = false;
      }
    },
  },
  beforeUnmount() {
    // Clear the interval when the component is destroyed
    clearInterval(this.interval);
  },
};
</script>

<style scoped>
.user-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 20px;
}

.user-card {
  background: #f9f9f9;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  max-width: 500px;
  min-width: 200px;
  width: 100%;
}

.user-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.user-card h3 {
  margin: 0;
  color: #333;
}

.user-card p {
  margin: 5px 0;
  color: #555;
}

.user-photo {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: 50%;
  margin-bottom: 10px;
}
</style>
