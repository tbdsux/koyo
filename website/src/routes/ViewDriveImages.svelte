<script lang="ts">
	import { onMount } from 'svelte';

	export let baseUrl: string;
	export let file: string;

	let x = file.split('Z-')[1].split('.');
	x.pop();
	const hostname = x.join('.');
	console.log(hostname);

	let imgBlobUrl: string = '';

	onMount(async () => {
		const r = await fetch(`${baseUrl}/drive/files/${file}`);
		const imageData = await r.blob();
		imgBlobUrl = URL.createObjectURL(imageData);
	});
</script>

{#if imgBlobUrl == ''}
	<li class="h-56 w-full bg-gray-100 animate-pulse" />
{:else}
	<li class="h-56 w-full relative bg-gray-100 group">
		<div
			class="bg-black/50 absolute w-full h-full hidden group-hover:flex items-center justify-center z-40"
		>
			<a
				href={`${baseUrl}/drive/files/${file}`}
				target="_blank"
				rel="noreferrer noopener"
				class="py-1 px-4 rounded-lg bg-gray-200 hover:bg-gray-300 duration-300">view</a
			>
		</div>

		<img src={imgBlobUrl} alt={file} class="absolute z-20 w-full h-full object-cover" />
		<span class="text-sm bg-gray-200 py-1 px-3 rounded-lg absolute z-30 top-2 left-2"
			>{hostname}</span
		>
	</li>
{/if}
