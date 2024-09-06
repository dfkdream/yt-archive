<script>
    import { Block, Navbar, Page, NavbarBackLink } from 'konsta/svelte';
</script>

<Page>
    <Navbar medium transparent title="Test">
        <NavbarBackLink slot="left" text="Back" onClick={()=>{history.back()}} />
    </Navbar>

    <Block strong inset>This is just a test page!</Block>
</Page>