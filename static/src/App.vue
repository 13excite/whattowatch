<template>

  <div id="app" class="center">
      <h1 class="title">{{ film.title }}</h1>
      <h2 class="genre">{{ film.genre }}</h2>
    <h2 class="rating">KP: {{ film.rating_kp[0] }} {{ film.rating_imdb[0] }}</h2>
      <div class="poster">
        <img style="width: auto; max-height: 450px;" :src=film.poster_link>
      </div>
      <div class="country">
        <h2 v-for="c in film.country " :key="c">Страна: {{ c }} </h2>
      </div>
    <div class="kp_link">
      <a target="_blank" rel="noopener noreferrer" :href=film.link_to_kp>
        <button class="btn btn-warning btn-xs btn-lg">Go to Kinopoisk</button>
      </a>

    </div>
    <div class="next">
      <button class="btn btn-xs btn-warning btn-lg" @click="update">Next</button>
    </div>

  </div>
</template>

<script>
  import axios from 'axios'

  export default {
    name: 'app',
    data () {
      return {
        film: {}
      }
    },
    created () {
      axios.get('http://127.0.0.1:8081/random').then(r => {
        this.film = r.data
      })
    },
    methods: {
      update () {
        axios.get('http://127.0.0.1:8081/random').then(r => {
          this.film = r.data
        })
      }
    }
  }
</script>

<style>
html body {
  background: rgb(0,0,0);
  background: linear-gradient(90deg, rgba(0,0,0,0.5606617647058824) 0%, rgba(9,121,112,1) 40%, rgba(0,212,255,1) 100%);

}
.center {
  text-align: center;
}
.country>h2 { display: inline }

.next {
    margin-top: 10px;
}
div.kp_link, .country {
    margin-top: 20px;
}
</style>
