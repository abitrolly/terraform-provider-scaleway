package scaleway

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nicolai86/scaleway-sdk"
)

func dataSourceScalewayVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScalewayVolumeRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of the volume",
			},
			"size_in_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the size of the volume in GB",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the type of backing storage",
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceScalewayVolumeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client).scaleway

	name := d.Get("name").(string)

	volumes, err := client.GetVolumes()
	if err != nil {
		if serr, ok := err.(api.APIError); ok {
			log.Printf("[DEBUG] Error obtaining Volumes: %q\n", serr.APIMessage)
		}

		return fmt.Errorf("Error obtaining Volumes: %+v", err)
	}

	var volume *api.Volume
	for _, v := range *volumes {
		if v.Name == name {
			volume = &v
			break
		}
	}

	if volume == nil {
		return fmt.Errorf("Couldn't locate a Volume with the name %q!", name)
	}

	d.SetId(volume.Identifier)

	d.Set("name", volume.Name)
	d.Set("size_in_gb", int(uint64(volume.Size)/gb))
	d.Set("type", volume.VolumeType)

	if volume.Server != nil {
		d.Set("server", volume.Server.Identifier)
	} else {
		d.Set("server", "")
	}

	return nil
}
