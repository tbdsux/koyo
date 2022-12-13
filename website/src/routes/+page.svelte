<script lang="ts">
	import { apiUrl, type APIResponse } from '$lib/api';
	import type { PageServerData } from './$types';
	import RenderImage from './RenderImage.svelte';

	export let data: PageServerData;

	let websiteUrl = '';
	let optionsHeight = 800;
	let optionsWidth = 1280;
	let optionsFullpage = false;
	let optionsImageType = 'png';
	let optionsDriver = 'playwright';

	let showOptions = false;
	let error = false;
	let errorData: APIResponse | null = null;
	let fetching = false;

	let imageData = '';

	const fetchScreenshot = async () => {
		if (!websiteUrl) return;

		imageData = '';
		fetching = true;

		if (optionsDriver === 'playwright') {
			optionsImageType = 'png';
		}

		const params = new URLSearchParams({
			height: optionsHeight.toString(),
			width: optionsWidth.toString(),
			imageType: optionsImageType,
			fullPage: optionsFullpage.toString(),
			driver: optionsDriver
		});
		const paramQueries = params.toString();
		const url = paramQueries != '' ? `/screenshot?${paramQueries}` : '/screenshot';

		const r = await fetch(apiUrl + url, {
			method: 'POST',
			headers: {
				'content-type': 'application/json'
			},
			body: JSON.stringify({ website: websiteUrl })
		});

		fetching = false;

		if (!r.ok) {
			error = true;
			errorData = await r.json();
			return;
		}

		error = false;
		errorData = null;
		const imgBlog = await r.blob();
		imageData = URL.createObjectURL(imgBlog);
	};
</script>

<header class="py-4 w-5/6 mx-auto text-center">
	<span class="inline-flex items-center">
		<img src="/favicon.png" alt="Koyo" class="object-fit h-12 w-12" />
		<h1 class="font-black text-xl text-blue-500">Koyo</h1>
	</span>
	<p class="text-neutral-600">Website screenshot service api</p>
</header>

<hr class="w-1/2 mx-auto my-3 border-gray-100" />

<div>
	<div class="flex flex-col text-left w-full sm:w-11/12 md:w-5/6 lg:w-4/5 xl:w-3/5 mx-auto">
		<label for="website" class="text-neutral-700">Enter website url:</label>
		<div class="flex items-center flex-col md:flex-row">
			<input
				bind:value={websiteUrl}
				type="text"
				id="website"
				name="website"
				class="m-1 py-2 px-5 rounded-xl border-2 outline-none hover:border-blue-500 focus:border-blue-500 w-full"
				placeholder="https://www.deta.sh"
			/>

			<button
				on:click={fetchScreenshot}
				disabled={fetching}
				class="m-1 inline-flex items-center text-sm bg-blue-400 hover:bg-blue-500 py-2 px-8 rounded-xl text-white duration-300 disabled:opacity-80 disabled:hover:bg-blue-400"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M9 17.25v1.007a3 3 0 01-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0115 18.257V17.25m6-12V15a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 15V5.25m18 0A2.25 2.25 0 0018.75 3H5.25A2.25 2.25 0 003 5.25m18 0V12a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 12V5.25"
					/>
				</svg>

				<span class="ml-1">
					{#if fetching}
						Fetching...
					{:else}
						Screenshot
					{/if}</span
				>
			</button>
		</div>

		<div class="mt-2">
			<button class="text-xs text-gray-600" on:click={() => (showOptions = !showOptions)}>
				{showOptions ? 'Hide' : 'Show'} options</button
			>

			{#if showOptions}
				<div
					class="flex items-center justify-center flex-wrap p-2 rounded-xl border border-gray-100"
				>
					<div class="inline-flex flex-col m-2">
						<label for="width" class="text-sm text-gray-700">Width</label>
						<input
							bind:value={optionsWidth}
							type="number"
							min={240}
							max={16834}
							name="width"
							id="width"
							class="py-2 px-3 text-sm border rounded-xl w-32 disabled:opacity-80"
						/>
					</div>

					<div class="inline-flex flex-col m-2">
						<label for="height" class="text-sm text-gray-700">Height</label>
						<input
							bind:value={optionsHeight}
							disabled={optionsFullpage}
							type="number"
							min={240}
							max={16834}
							name="height"
							id="height"
							class="py-2 px-3 text-sm border rounded-xl w-32"
						/>
					</div>

					<div class="inline-flex flex-col m-2">
						<label for="width" class="text-sm text-gray-700">Image Type</label>
						<select
							bind:value={optionsImageType}
							class="py-2 px-3 text-sm border rounded-xl w-32 disabled:opacity-80 bg-white"
						>
							<option value="png">png</option>
							<option value="jpeg">jpeg</option>
							<option disabled={optionsDriver === 'playwright'} value="webp">webp</option>
						</select>
					</div>

					<div class="inline-flex flex-col m-2">
						<label for="width" class="text-sm text-gray-700">Driver</label>
						<select
							bind:value={optionsDriver}
							class="py-2 px-3 text-sm border rounded-xl w-32 disabled:opacity-80 bg-white"
						>
							<option value="playwright">playwright</option>
							<option value="puppeteer">puppeteer</option>
						</select>
					</div>

					<div class="inline-flex items-center flex-row m-2">
						<input
							bind:checked={optionsFullpage}
							type="checkbox"
							name="fullpage"
							id="fullpage"
							class="h-4 w-4 text-sm border rounded-xl"
						/>
						<label for="fullpage" class="text-sm ml-1 text-gray-700">Full Page</label>
					</div>
				</div>
			{/if}
		</div>

		<pre
			class="mt-4 overflow-auto text-left text-sm bg-gray-100 text-gray-700 py-2 px-3 rounded-xl">
{`curl -X POST \\
	'${data.baseUrl}/screenshot?height=${optionsHeight}&width=${optionsWidth}&imageType=${optionsImageType}&fullPage=${optionsFullpage}' \\
	--header 'Content-Type: application/json' \\
	--data-raw '{"website":"${websiteUrl}"}'`}
		</pre>
	</div>

	<div class="mt-12 h-screen w-full relative overflow-auto bg-gray-50 mb-20">
		{#if fetching}
			<div class="h-full w-full animate-pulse bg-gray-300" />
		{:else if error}
			<p>{JSON.stringify(errorData)}</p>
		{:else if imageData != ''}
			<RenderImage bind:imageData bind:websiteUrl />
		{/if}
	</div>
</div>
