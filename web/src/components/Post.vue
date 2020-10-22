<template>
  <div>
    <h1 class="title">{{ post.Title }}</h1>
    <span>{{ post.Source }}</span> |
    <a :href="post.Url" target="_blank">{{ $t("post.link") }}</a>
    <b-card v-if="post.Content!==''" class="postcontext" style="margin: 15px" v-html="post.Content">
    </b-card>
    <b-card v-else class="postcontext" style="margin: 15px" v-html="post.Description">
    </b-card>
  </div>
</template>

<script>
export default {
  name: "Post",
  props: ['post', 'nextpost'],
}
</script>

<style scoped>

</style>