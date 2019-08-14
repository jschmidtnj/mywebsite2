<template>
  <div class="main-footer">
    <footer>
      <b-container>
        <p>&copy; Joshua Schmidt 2019</p>
        <ul class="list-inline">
          <li class="list-inline-item">
            <a href="/privacy">Privacy</a>
          </li>
          <li class="list-inline-item">
            <a :href="githuburl">Github</a>
          </li>
        </ul>
      </b-container>
    </footer>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
const githuburl = process.env.githuburl
export default Vue.extend({
  name: 'Footer',
  data() {
    return {
      githuburl: githuburl
    }
  }
})
</script>

<style lang="scss">
.main-footer {
  footer {
    width: 100%;
    padding: 25px;
    text-align: center;
    color: rgba(255, 255, 255, 0.3);
    background-color: #222222;
  }

  footer a {
    color: #fdcc52;
    -webkit-transition: all 0.35s;
    -moz-transition: all 0.35s;
    transition: all 0.35s;
  }

  footer a:hover,
  footer a:focus {
    color: #fcbd20;
  }

  footer p {
    font-size: 12px;
    margin: 0;
  }

  footer ul {
    margin-bottom: 0;
  }

  footer ul li a {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.3);
  }

  footer ul li a:hover,
  footer ul li a:focus,
  footer ul li a:active,
  footer ul li a.active {
    text-decoration: none;
  }
}
</style>
