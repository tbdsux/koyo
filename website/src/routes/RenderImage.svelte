<script lang="ts">
	export let imageData: string;
	export let websiteUrl: string;
</script>

<img src={imageData} alt={websiteUrl} class="shadow-xl border mx-auto" />
