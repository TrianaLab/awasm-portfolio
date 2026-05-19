<script lang="ts">
  type Mode = 'terminal' | 'resume';
  let { mode, onChange }: { mode: Mode; onChange: (m: Mode) => void } = $props();
</script>

<div class="segmented" role="tablist" aria-label="View mode">
  <button
    role="tab"
    aria-selected={mode === 'terminal'}
    class:active={mode === 'terminal'}
    onclick={() => onChange('terminal')}
  >
    <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor" aria-hidden="true">
      <path d="M4 5a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h16a1 1 0 0 0 1-1V6a1 1 0 0 0-1-1H4Zm2 4 3 3-3 3m4 0h6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
    <span>Terminal</span>
  </button>
  <button
    role="tab"
    aria-selected={mode === 'resume'}
    class:active={mode === 'resume'}
    onclick={() => onChange('resume')}
  >
    <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor" aria-hidden="true">
      <path d="M6 2h9l5 5v15a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V3a1 1 0 0 1 1-1Zm9 0v5h5M8 13h8M8 17h8M8 9h4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
    <span>Resume</span>
  </button>
</div>

<style>
  .segmented {
    display: inline-flex;
    background: var(--color-bg-subtle);
    border-radius: var(--radius-sm);
    padding: 2px;
    border: 1px solid var(--color-border);
  }

  .segmented button {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0 0.65rem;
    color: var(--color-text-muted);
    font-size: 0.82rem;
    font-weight: 500;
    border-radius: 4px;
    border: none;
    background: transparent;
    cursor: pointer;
    transition: all var(--transition);
  }

  .segmented button:hover {
    color: var(--color-text);
  }

  .segmented button.active {
    background: var(--color-bg-elevated);
    color: var(--color-text);
    box-shadow: var(--shadow-md);
  }

  @media (max-width: 480px) {
    .segmented button span {
      display: none;
    }
  }
</style>
