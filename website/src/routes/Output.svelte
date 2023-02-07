<script lang="ts">
	export let imageUrl: string;
	export let websiteUrl: string;

	export let optionsNoOutput: boolean;
	const optNoOutput = optionsNoOutput; // we re-assign the value so that it won't re-render if `optionsNoOutput` changes value
</script>

{#if optNoOutput}
	<span class="text-sm p-4 text-center">Successfully saved screenshot to Drive.</span>
{:else}
	<img src={imageUrl} alt={websiteUrl} class="shadow-xl border mx-auto" />
{/if}
