<script>
    import { Block, BlockTitle, Navbar, Page, List, ListItem, Button, Dialog, DialogButton } from 'konsta/svelte';

    let alertString = "";
    let alertOpened = false;
    let confirmOpened = false;
</script>

<Page>
    <Navbar medium transparent title="SvelteKit PWA" />

    <Block strong inset>Visit <a href="https://github.com/dfkdream/sveltekit-pwa-template">https://github.com/dfkdream/sveltekit-pwa-template</a> for source code</Block>

    <BlockTitle>Block Title</BlockTitle>
    <List strong inset>
        <ListItem href="https://blog.dfkdream.dev" title="Blog" />
        <ListItem href="https://about.dfkdream.dev" title="About" />
        <ListItem href="/test" title="Test" />
    </List>

    <Block strong inset class="flex space-x-4">
        <Button onClick={()=>{alertString="Hello, world!"; alertOpened=true}}>Alert</Button>
        <Button onClick={()=>{confirmOpened=true}}>Confirm</Button>
    </Block>

    <BlockTitle>Here are some lorem ipsum</BlockTitle>
    <Block strong inset>
        Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quia, quo rem beatae, delectus eligendi est saepe molestias perferendis suscipit, commodi labore ipsa non quasi eum magnam neque ducimus! Quasi, numquam.

        Maiores culpa, itaque! Eaque natus ab cum ipsam numquam blanditiis a, quia, molestiae aut laudantium recusandae ipsa. Ad iste ex asperiores ipsa, mollitia perferendis consectetur quam eaque, voluptate laboriosam unde.

        Sed odit quis aperiam temporibus vitae necessitatibus, laboriosam, exercitationem dolores odio sapiente provident. Accusantium id, itaque aliquam libero ipsum eos fugiat distinctio laboriosam exercitationem sequi facere quas quidem magnam reprehenderit.

        Pariatur corporis illo, amet doloremque. Ab veritatis sunt nisi consectetur error modi, nam illo et nostrum quia aliquam ipsam vitae facere voluptates atque similique odit mollitia, rerum placeat nobis est.

        Et impedit soluta minus a autem adipisci cupiditate eius dignissimos nihil officia dolore voluptatibus aperiam reprehenderit esse facilis labore qui, officiis consectetur. Ipsa obcaecati aspernatur odio assumenda veniam, ipsum alias.

        Lorem ipsum dolor sit amet, consectetur adipisicing elit. Culpa ipsa debitis sed nihil eaque dolore cum iste quibusdam, accusamus doloribus, tempora quia quos voluptatibus corporis officia at quas dolorem earum!

        Quod soluta eos inventore magnam suscipit enim at hic in maiores temporibus pariatur tempora minima blanditiis vero autem est perspiciatis totam dolorum, itaque repellat? Nobis necessitatibus aut odit aliquam adipisci.

        Tenetur delectus perspiciatis ex numquam, unde corrupti velit! Quam aperiam, animi fuga veritatis consectetur, voluptatibus atque consequuntur dignissimos itaque, sint impedit cum cumque at. Adipisci sint, iusto blanditiis ullam? Vel?

        Dignissimos velit officia quibusdam! Eveniet beatae, aut, omnis temporibus consequatur expedita eaque aliquid quos accusamus fugiat id iusto autem obcaecati repellat fugit cupiditate suscipit natus quas doloribus? Temporibus necessitatibus, libero.

        Architecto quisquam ipsa fugit facere, repudiandae asperiores vitae obcaecati possimus, labore excepturi reprehenderit consectetur perferendis, ullam quidem hic, repellat fugiat eaque fuga. Consectetur in eveniet, deleniti recusandae omnis eum quas?

        Quos nulla consequatur quo, officia quaerat. Nulla voluptatum, assumenda quibusdam, placeat cum aut illo deleniti dolores commodi odio ipsam, recusandae est pariatur veniam repudiandae blanditiis. Voluptas unde deleniti quisquam, nobis?

        Atque qui quaerat quasi officia molestiae, molestias totam incidunt reprehenderit laboriosam facilis veritatis, non iusto! Dolore ipsam obcaecati voluptates minima maxime minus qui mollitia facere. Nostrum esse recusandae voluptatibus eligendi.
    </Block>

    <Dialog opened={alertOpened} onBackdropClick={()=>{alertOpened=false}}>
        <svelte:fragment slot="title">Alert</svelte:fragment>
        {alertString}
        <svelte:fragment slot="buttons">
            <DialogButton onClick={()=>{alertOpened=false}}>OK</DialogButton>
        </svelte:fragment>
    </Dialog>

    <Dialog opened={confirmOpened}>
        <svelte:fragment slot="title">Confirm</svelte:fragment>
        This button does nothing. Alright?
        <svelte:fragment slot="buttons">
            <DialogButton onClick={()=>{confirmOpened=false; alertString="Too bad!"; alertOpened=true}}>NO</DialogButton>
            <DialogButton onClick={()=>{confirmOpened=false}}>OK</DialogButton>
        </svelte:fragment>
    </Dialog>
</Page>
