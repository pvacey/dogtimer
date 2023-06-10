<script>
  import { onMount } from 'svelte';
  let proto = import.meta.env.VITE_PROTO
  let hostname = import.meta.env.VITE_HOSTNAME;
  let port = import.meta.env.VITE_PORT;
  let base_url = `${proto}://${hostname}:${port}`
  let previousLast = new Date().getTime();
  let lastTime = new Date().getTime();
  let hours = '';
  let minutes = '';
  let diff = 0;
  let height = window.innerHeight/3;
  export let name;

  const updateClock = () => {
    setTimeout(() => {
      let now = new Date().getTime();
      diff = (now - lastTime)/1000;
      hours   = Math.floor(diff/60/60);
      minutes = Math.floor(diff/60%60);

      if (now % 3 === 0) {
        getLatestTime()
      }
      height = window.innerHeight/3;
      updateClock()
    }, 1000);
  }

  const undo = () => {
    lastTime = previousLast;
    setLatestTime();
  }

  const reset = () => {
    if (diff > 300) {
      hours = 0
      minutes = 0
      diff = 0
      previousLast = lastTime;
      lastTime = new Date().getTime();
      setLatestTime();
    }
  }

  async function getLatestTime() {
    const resp1 = await fetch(`${base_url}/time/${name}`)
      resp1.json().then(data => {
        lastTime = data.Time;
    });
    const resp2  = await fetch(`${base_url}/time/${name}_previous`)
      resp2.json().then(data => {
        previousLast = data.Time;
    });
  } 

  async function setLatestTime() {
    await fetch(`${base_url}/time/${name}/${lastTime}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        }
    })
    await fetch(`${base_url}/time/${name}_previous/${previousLast}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        }
    })
  } 
    
  onMount(() => {
    updateClock();
    getLatestTime();

  });
</script>

<button class="card {name}" style="height:{height*.78}px" on:click={reset}>
  {#if hours === ''}
    loading...
  {:else if hours<1}
    {minutes} mins
  {:else if hours === 1}
    {hours} hour<br>{minutes} mins
  {:else}
    {hours} hours<br>{minutes} mins
  {/if}
</button>
<button class="undo" style="height:{height*.2}px" on:click={undo}>
  undo
</button>
