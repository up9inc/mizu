import React, { useCallback, useEffect, useState } from "react";
import { Backdrop, Box, Button, debounce, Fade, Modal } from "@mui/material";
import styles from "./TrafficStatsModal.module.sass";
import closeIcon from "assets/close.svg";
import { TrafficPieChart } from "./TrafficPieChart/TrafficPieChart";
import { TimelineBarChart } from "./TimelineBarChart/TimelineBarChart";
import spinnerImg from "assets/spinner.svg";
import refreshIcon from "assets/refresh.svg";
import { useCommonStyles } from "../../../helpers/commonStyle";

const modalStyle = {
  position: 'absolute',
  top: '6%',
  left: '50%',
  transform: 'translate(-50%, 0%)',
  width: '50vw',
  height: '82vh',
  bgcolor: 'background.paper',
  borderRadius: '5px',
  boxShadow: 24,
  p: 4,
  color: '#000',
};

export enum StatsMode {
  REQUESTS = "entriesCount",
  VOLUME = "volumeSizeBytes"
}

interface TrafficStatsModalProps {
  isOpen: boolean;
  onClose: () => void;
  getPieStatsDataApi: () => Promise<any>
  getTimelineStatsDataApi: () => Promise<any>
}

export const dataMock = {
  "pie": [
    {
      "name": "REDIS",
      "entriesCount": 101068,
      "volumeSizeBytes": 54646626,
      "color": "#a41e11",
      "methods": [
        {
          "name": "SET",
          "entriesCount": 33712,
          "volumeSizeBytes": 18349129
        },
        {
          "name": "GET",
          "entriesCount": 67356,
          "volumeSizeBytes": 36297497
        }
      ]
    },
    {
      "name": "gRPC",
      "entriesCount": 78262,
      "volumeSizeBytes": 352670423,
      "color": "#244c5a",
      "methods": [
        {
          "name": "POST",
          "entriesCount": 78262,
          "volumeSizeBytes": 352670423
        }
      ]
    },
    {
      "name": "HTTP",
      "entriesCount": 138007,
      "volumeSizeBytes": 627650751,
      "color": "#244c5a",
      "methods": [
        {
          "name": "GET",
          "entriesCount": 137988,
          "volumeSizeBytes": 627594425
        },
        {
          "name": "POST",
          "entriesCount": 18,
          "volumeSizeBytes": 54801
        },
        {
          "name": "HEAD",
          "entriesCount": 1,
          "volumeSizeBytes": 1525
        }
      ]
    },
    {
      "name": "KAFKA",
      "entriesCount": 60174,
      "volumeSizeBytes": 89401119,
      "color": "#000000",
      "methods": [
        {
          "name": "Metadata",
          "entriesCount": 19270,
          "volumeSizeBytes": 23214235
        },
        {
          "name": "Produce",
          "entriesCount": 5576,
          "volumeSizeBytes": 10938489
        },
        {
          "name": "CreateTopics",
          "entriesCount": 7893,
          "volumeSizeBytes": 7619820
        },
        {
          "name": "ListOffsets",
          "entriesCount": 358,
          "volumeSizeBytes": 329634
        },
        {
          "name": "ApiVersions",
          "entriesCount": 27077,
          "volumeSizeBytes": 47298941
        }
      ]
    },
    {
      "name": "GQL",
      "entriesCount": 173285,
      "volumeSizeBytes": 591359689,
      "color": "#e10098",
      "methods": [
        {
          "name": "POST",
          "entriesCount": 173285,
          "volumeSizeBytes": 591359689
        }
      ]
    },
    {
      "name": "AMQP",
      "entriesCount": 93029,
      "volumeSizeBytes": 66571843,
      "color": "#ff6600",
      "methods": [
        {
          "name": "queue declare",
          "entriesCount": 8255,
          "volumeSizeBytes": 4482281
        },
        {
          "name": "exchange declare",
          "entriesCount": 8255,
          "volumeSizeBytes": 4738207
        },
        {
          "name": "queue bind",
          "entriesCount": 8255,
          "volumeSizeBytes": 4300571
        },
        {
          "name": "basic consume",
          "entriesCount": 8255,
          "volumeSizeBytes": 4754630
        },
        {
          "name": "basic publish",
          "entriesCount": 16510,
          "volumeSizeBytes": 13059068
        },
        {
          "name": "connection close",
          "entriesCount": 8255,
          "volumeSizeBytes": 4110814
        },
        {
          "name": "connection start",
          "entriesCount": 11748,
          "volumeSizeBytes": 11863519
        },
        {
          "name": "basic deliver",
          "entriesCount": 23496,
          "volumeSizeBytes": 19262753
        }
      ]
    }
  ],
  "timeline": [
    {
      "protocols": [
        {
          "name": "AMQP",
          "entriesCount": 12943,
          "volumeSizeBytes": 9247792,
          "color": "#ff6600",
          "methods": [
            {
              "name": "queue bind",
              "entriesCount": 1156,
              "volumeSizeBytes": 602200
            },
            {
              "name": "basic consume",
              "entriesCount": 1156,
              "volumeSizeBytes": 665802
            },
            {
              "name": "basic publish",
              "entriesCount": 2312,
              "volumeSizeBytes": 1828643
            },
            {
              "name": "connection start",
              "entriesCount": 1617,
              "volumeSizeBytes": 1632929
            },
            {
              "name": "connection close",
              "entriesCount": 1156,
              "volumeSizeBytes": 575637
            },
            {
              "name": "basic deliver",
              "entriesCount": 3234,
              "volumeSizeBytes": 2651438
            },
            {
              "name": "queue declare",
              "entriesCount": 1156,
              "volumeSizeBytes": 627646
            },
            {
              "name": "exchange declare",
              "entriesCount": 1156,
              "volumeSizeBytes": 663497
            }
          ]
        },
        {
          "name": "HTTP",
          "entriesCount": 19196,
          "volumeSizeBytes": 87289157,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 2,
              "volumeSizeBytes": 5891
            },
            {
              "name": "GET",
              "entriesCount": 19194,
              "volumeSizeBytes": 87283266
            }
          ]
        },
        {
          "name": "REDIS",
          "entriesCount": 14028,
          "volumeSizeBytes": 7584244,
          "color": "#a41e11",
          "methods": [
            {
              "name": "SET",
              "entriesCount": 4675,
              "volumeSizeBytes": 2544085
            },
            {
              "name": "GET",
              "entriesCount": 9353,
              "volumeSizeBytes": 5040159
            }
          ]
        },
        {
          "name": "GQL",
          "entriesCount": 23836,
          "volumeSizeBytes": 81340349,
          "color": "#e10098",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 23836,
              "volumeSizeBytes": 81340349
            }
          ]
        },
        {
          "name": "KAFKA",
          "entriesCount": 8343,
          "volumeSizeBytes": 12438497,
          "color": "#000000",
          "methods": [
            {
              "name": "ApiVersions",
              "entriesCount": 3783,
              "volumeSizeBytes": 6580302
            },
            {
              "name": "CreateTopics",
              "entriesCount": 1104,
              "volumeSizeBytes": 1065840
            },
            {
              "name": "Produce",
              "entriesCount": 756,
              "volumeSizeBytes": 1526458
            },
            {
              "name": "Metadata",
              "entriesCount": 2700,
              "volumeSizeBytes": 3265897
            }
          ]
        },
        {
          "name": "gRPC",
          "entriesCount": 10630,
          "volumeSizeBytes": 47921775,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 10630,
              "volumeSizeBytes": 47921775
            }
          ]
        }
      ],
      "timestamp": 1656518400000
    },
    {
      "protocols": [
        {
          "name": "gRPC",
          "entriesCount": 3643,
          "volumeSizeBytes": 16422934,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 3643,
              "volumeSizeBytes": 16422934
            }
          ]
        },
        {
          "name": "AMQP",
          "entriesCount": 4502,
          "volumeSizeBytes": 3205857,
          "color": "#ff6600",
          "methods": [
            {
              "name": "queue bind",
              "entriesCount": 408,
              "volumeSizeBytes": 212540
            },
            {
              "name": "connection close",
              "entriesCount": 407,
              "volumeSizeBytes": 202662
            },
            {
              "name": "basic deliver",
              "entriesCount": 1098,
              "volumeSizeBytes": 900166
            },
            {
              "name": "basic consume",
              "entriesCount": 408,
              "volumeSizeBytes": 234976
            },
            {
              "name": "basic publish",
              "entriesCount": 816,
              "volumeSizeBytes": 645415
            },
            {
              "name": "connection start",
              "entriesCount": 549,
              "volumeSizeBytes": 554410
            },
            {
              "name": "queue declare",
              "entriesCount": 408,
              "volumeSizeBytes": 221514
            },
            {
              "name": "exchange declare",
              "entriesCount": 408,
              "volumeSizeBytes": 234174
            }
          ]
        },
        {
          "name": "REDIS",
          "entriesCount": 4844,
          "volumeSizeBytes": 2618917,
          "color": "#a41e11",
          "methods": [
            {
              "name": "SET",
              "entriesCount": 1613,
              "volumeSizeBytes": 877738
            },
            {
              "name": "GET",
              "entriesCount": 3231,
              "volumeSizeBytes": 1741179
            }
          ]
        },
        {
          "name": "KAFKA",
          "entriesCount": 2878,
          "volumeSizeBytes": 4252557,
          "color": "#000000",
          "methods": [
            {
              "name": "CreateTopics",
              "entriesCount": 382,
              "volumeSizeBytes": 369996
            },
            {
              "name": "Metadata",
              "entriesCount": 947,
              "volumeSizeBytes": 1178186
            },
            {
              "name": "ApiVersions",
              "entriesCount": 1278,
              "volumeSizeBytes": 2235971
            },
            {
              "name": "Produce",
              "entriesCount": 271,
              "volumeSizeBytes": 468404
            }
          ]
        },
        {
          "name": "GQL",
          "entriesCount": 8168,
          "volumeSizeBytes": 27875694,
          "color": "#e10098",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 8168,
              "volumeSizeBytes": 27875694
            }
          ]
        },
        {
          "name": "HTTP",
          "entriesCount": 6626,
          "volumeSizeBytes": 30138355,
          "color": "#244c5a",
          "methods": [
            {
              "name": "GET",
              "entriesCount": 6625,
              "volumeSizeBytes": 30135398
            },
            {
              "name": "POST",
              "entriesCount": 1,
              "volumeSizeBytes": 2957
            }
          ]
        }
      ],
      "timestamp": 1656511200000
    },
    {
      "protocols": [
        {
          "name": "AMQP",
          "entriesCount": 10979,
          "volumeSizeBytes": 7852775,
          "color": "#ff6600",
          "methods": [
            {
              "name": "queue bind",
              "entriesCount": 976,
              "volumeSizeBytes": 508428
            },
            {
              "name": "basic consume",
              "entriesCount": 976,
              "volumeSizeBytes": 562128
            },
            {
              "name": "basic publish",
              "entriesCount": 1952,
              "volumeSizeBytes": 1543956
            },
            {
              "name": "connection close",
              "entriesCount": 977,
              "volumeSizeBytes": 486508
            },
            {
              "name": "connection start",
              "entriesCount": 1382,
              "volumeSizeBytes": 1395601
            },
            {
              "name": "basic deliver",
              "entriesCount": 2764,
              "volumeSizeBytes": 2266027
            },
            {
              "name": "queue declare",
              "entriesCount": 976,
              "volumeSizeBytes": 529929
            },
            {
              "name": "exchange declare",
              "entriesCount": 976,
              "volumeSizeBytes": 560198
            }
          ]
        },
        {
          "name": "KAFKA",
          "entriesCount": 7059,
          "volumeSizeBytes": 10308937,
          "color": "#000000",
          "methods": [
            {
              "name": "ApiVersions",
              "entriesCount": 3174,
              "volumeSizeBytes": 5555838
            },
            {
              "name": "Produce",
              "entriesCount": 656,
              "volumeSizeBytes": 1124683
            },
            {
              "name": "Metadata",
              "entriesCount": 2246,
              "volumeSizeBytes": 2684175
            },
            {
              "name": "CreateTopics",
              "entriesCount": 907,
              "volumeSizeBytes": 874273
            },
            {
              "name": "ListOffsets",
              "entriesCount": 76,
              "volumeSizeBytes": 69968
            }
          ]
        },
        {
          "name": "gRPC",
          "entriesCount": 9379,
          "volumeSizeBytes": 42268574,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 9379,
              "volumeSizeBytes": 42268574
            }
          ]
        },
        {
          "name": "REDIS",
          "entriesCount": 11904,
          "volumeSizeBytes": 6436740,
          "color": "#a41e11",
          "methods": [
            {
              "name": "SET",
              "entriesCount": 3973,
              "volumeSizeBytes": 2162784
            },
            {
              "name": "GET",
              "entriesCount": 7931,
              "volumeSizeBytes": 4273956
            }
          ]
        },
        {
          "name": "HTTP",
          "entriesCount": 16212,
          "volumeSizeBytes": 73761919,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 4,
              "volumeSizeBytes": 11925
            },
            {
              "name": "GET",
              "entriesCount": 16208,
              "volumeSizeBytes": 73749994
            }
          ]
        },
        {
          "name": "GQL",
          "entriesCount": 20552,
          "volumeSizeBytes": 70142018,
          "color": "#e10098",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 20552,
              "volumeSizeBytes": 70142018
            }
          ]
        }
      ],
      "timestamp": 1656561600000
    },
    {
      "protocols": [
        {
          "name": "HTTP",
          "entriesCount": 19199,
          "volumeSizeBytes": 87345193,
          "color": "#244c5a",
          "methods": [
            {
              "name": "GET",
              "entriesCount": 19196,
              "volumeSizeBytes": 87337770
            },
            {
              "name": "POST",
              "entriesCount": 2,
              "volumeSizeBytes": 5898
            },
            {
              "name": "HEAD",
              "entriesCount": 1,
              "volumeSizeBytes": 1525
            }
          ]
        },
        {
          "name": "KAFKA",
          "entriesCount": 8213,
          "volumeSizeBytes": 12030942,
          "color": "#000000",
          "methods": [
            {
              "name": "ApiVersions",
              "entriesCount": 3651,
              "volumeSizeBytes": 6434841
            },
            {
              "name": "Produce",
              "entriesCount": 769,
              "volumeSizeBytes": 1325530
            },
            {
              "name": "Metadata",
              "entriesCount": 2598,
              "volumeSizeBytes": 3125386
            },
            {
              "name": "CreateTopics",
              "entriesCount": 1069,
              "volumeSizeBytes": 1029156
            },
            {
              "name": "ListOffsets",
              "entriesCount": 126,
              "volumeSizeBytes": 116029
            }
          ]
        },
        {
          "name": "gRPC",
          "entriesCount": 11058,
          "volumeSizeBytes": 49834224,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 11058,
              "volumeSizeBytes": 49834224
            }
          ]
        },
        {
          "name": "REDIS",
          "entriesCount": 14095,
          "volumeSizeBytes": 7622590,
          "color": "#a41e11",
          "methods": [
            {
              "name": "GET",
              "entriesCount": 9385,
              "volumeSizeBytes": 5057500
            },
            {
              "name": "SET",
              "entriesCount": 4710,
              "volumeSizeBytes": 2565090
            }
          ]
        },
        {
          "name": "GQL",
          "entriesCount": 24363,
          "volumeSizeBytes": 83126241,
          "color": "#e10098",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 24363,
              "volumeSizeBytes": 83126241
            }
          ]
        },
        {
          "name": "AMQP",
          "entriesCount": 13030,
          "volumeSizeBytes": 9301829,
          "color": "#ff6600",
          "methods": [
            {
              "name": "basic deliver",
              "entriesCount": 3236,
              "volumeSizeBytes": 2652966
            },
            {
              "name": "connection close",
              "entriesCount": 1168,
              "volumeSizeBytes": 581597
            },
            {
              "name": "queue declare",
              "entriesCount": 1168,
              "volumeSizeBytes": 634157
            },
            {
              "name": "exchange declare",
              "entriesCount": 1168,
              "volumeSizeBytes": 670361
            },
            {
              "name": "queue bind",
              "entriesCount": 1168,
              "volumeSizeBytes": 608434
            },
            {
              "name": "basic consume",
              "entriesCount": 1168,
              "volumeSizeBytes": 672702
            },
            {
              "name": "basic publish",
              "entriesCount": 2336,
              "volumeSizeBytes": 1847668
            },
            {
              "name": "connection start",
              "entriesCount": 1618,
              "volumeSizeBytes": 1633944
            }
          ]
        }
      ],
      "timestamp": 1656554400000
    },
    {
      "protocols": [
        {
          "name": "REDIS",
          "entriesCount": 14094,
          "volumeSizeBytes": 7620277,
          "color": "#a41e11",
          "methods": [
            {
              "name": "SET",
              "entriesCount": 4701,
              "volumeSizeBytes": 2558643
            },
            {
              "name": "GET",
              "entriesCount": 9393,
              "volumeSizeBytes": 5061634
            }
          ]
        },
        {
          "name": "gRPC",
          "entriesCount": 11201,
          "volumeSizeBytes": 50461853,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 11201,
              "volumeSizeBytes": 50461853
            }
          ]
        },
        {
          "name": "HTTP",
          "entriesCount": 19191,
          "volumeSizeBytes": 87243765,
          "color": "#244c5a",
          "methods": [
            {
              "name": "GET",
              "entriesCount": 19188,
              "volumeSizeBytes": 87234720
            },
            {
              "name": "POST",
              "entriesCount": 3,
              "volumeSizeBytes": 9045
            }
          ]
        },
        {
          "name": "KAFKA",
          "entriesCount": 8645,
          "volumeSizeBytes": 12654533,
          "color": "#000000",
          "methods": [
            {
              "name": "ApiVersions",
              "entriesCount": 3931,
              "volumeSizeBytes": 6884479
            },
            {
              "name": "Metadata",
              "entriesCount": 2695,
              "volumeSizeBytes": 3213949
            },
            {
              "name": "CreateTopics",
              "entriesCount": 1104,
              "volumeSizeBytes": 1066678
            },
            {
              "name": "Produce",
              "entriesCount": 815,
              "volumeSizeBytes": 1397354
            },
            {
              "name": "ListOffsets",
              "entriesCount": 100,
              "volumeSizeBytes": 92073
            }
          ]
        },
        {
          "name": "AMQP",
          "entriesCount": 13040,
          "volumeSizeBytes": 9326096,
          "color": "#ff6600",
          "methods": [
            {
              "name": "exchange declare",
              "entriesCount": 1160,
              "volumeSizeBytes": 665820
            },
            {
              "name": "connection close",
              "entriesCount": 1160,
              "volumeSizeBytes": 577675
            },
            {
              "name": "queue bind",
              "entriesCount": 1160,
              "volumeSizeBytes": 604353
            },
            {
              "name": "basic consume",
              "entriesCount": 1160,
              "volumeSizeBytes": 668138
            },
            {
              "name": "basic publish",
              "entriesCount": 2320,
              "volumeSizeBytes": 1835086
            },
            {
              "name": "connection start",
              "entriesCount": 1640,
              "volumeSizeBytes": 1656109
            },
            {
              "name": "basic deliver",
              "entriesCount": 3280,
              "volumeSizeBytes": 2689035
            },
            {
              "name": "queue declare",
              "entriesCount": 1160,
              "volumeSizeBytes": 629880
            }
          ]
        },
        {
          "name": "GQL",
          "entriesCount": 24435,
          "volumeSizeBytes": 83399318,
          "color": "#e10098",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 24435,
              "volumeSizeBytes": 83399318
            }
          ]
        }
      ],
      "timestamp": 1656547200000
    },
    {
      "protocols": [
        {
          "name": "KAFKA",
          "entriesCount": 8279,
          "volumeSizeBytes": 12322152,
          "color": "#000000",
          "methods": [
            {
              "name": "CreateTopics",
              "entriesCount": 1085,
              "volumeSizeBytes": 1050137
            },
            {
              "name": "ListOffsets",
              "entriesCount": 56,
              "volumeSizeBytes": 51564
            },
            {
              "name": "ApiVersions",
              "entriesCount": 3747,
              "volumeSizeBytes": 6536884
            },
            {
              "name": "Produce",
              "entriesCount": 754,
              "volumeSizeBytes": 1520023
            },
            {
              "name": "Metadata",
              "entriesCount": 2637,
              "volumeSizeBytes": 3163544
            }
          ]
        },
        {
          "name": "HTTP",
          "entriesCount": 19203,
          "volumeSizeBytes": 87316591,
          "color": "#244c5a",
          "methods": [
            {
              "name": "GET",
              "entriesCount": 19199,
              "volumeSizeBytes": 87303406
            },
            {
              "name": "POST",
              "entriesCount": 4,
              "volumeSizeBytes": 13185
            }
          ]
        },
        {
          "name": "gRPC",
          "entriesCount": 10865,
          "volumeSizeBytes": 48951151,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 10865,
              "volumeSizeBytes": 48951151
            }
          ]
        },
        {
          "name": "AMQP",
          "entriesCount": 12806,
          "volumeSizeBytes": 9215819,
          "color": "#ff6600",
          "methods": [
            {
              "name": "basic publish",
              "entriesCount": 2218,
              "volumeSizeBytes": 1754445
            },
            {
              "name": "connection start",
              "entriesCount": 1681,
              "volumeSizeBytes": 1697495
            },
            {
              "name": "basic deliver",
              "entriesCount": 3362,
              "volumeSizeBytes": 2756283
            },
            {
              "name": "queue declare",
              "entriesCount": 1109,
              "volumeSizeBytes": 602181
            },
            {
              "name": "exchange declare",
              "entriesCount": 1109,
              "volumeSizeBytes": 636561
            },
            {
              "name": "connection close",
              "entriesCount": 1109,
              "volumeSizeBytes": 552290
            },
            {
              "name": "queue bind",
              "entriesCount": 1109,
              "volumeSizeBytes": 577771
            },
            {
              "name": "basic consume",
              "entriesCount": 1109,
              "volumeSizeBytes": 638793
            }
          ]
        },
        {
          "name": "GQL",
          "entriesCount": 24047,
          "volumeSizeBytes": 82066023,
          "color": "#e10098",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 24047,
              "volumeSizeBytes": 82066023
            }
          ]
        },
        {
          "name": "REDIS",
          "entriesCount": 14070,
          "volumeSizeBytes": 7607449,
          "color": "#a41e11",
          "methods": [
            {
              "name": "SET",
              "entriesCount": 4694,
              "volumeSizeBytes": 2554913
            },
            {
              "name": "GET",
              "entriesCount": 9376,
              "volumeSizeBytes": 5052536
            }
          ]
        }
      ],
      "timestamp": 1656540000000
    },
    {
      "protocols": [
        {
          "name": "AMQP",
          "entriesCount": 12964,
          "volumeSizeBytes": 9254988,
          "color": "#ff6600",
          "methods": [
            {
              "name": "connection start",
              "entriesCount": 1610,
              "volumeSizeBytes": 1625873
            },
            {
              "name": "queue declare",
              "entriesCount": 1162,
              "volumeSizeBytes": 630900
            },
            {
              "name": "connection close",
              "entriesCount": 1162,
              "volumeSizeBytes": 578608
            },
            {
              "name": "basic consume",
              "entriesCount": 1162,
              "volumeSizeBytes": 669239
            },
            {
              "name": "basic publish",
              "entriesCount": 2324,
              "volumeSizeBytes": 1838199
            },
            {
              "name": "basic deliver",
              "entriesCount": 3220,
              "volumeSizeBytes": 2639891
            },
            {
              "name": "exchange declare",
              "entriesCount": 1162,
              "volumeSizeBytes": 666935
            },
            {
              "name": "queue bind",
              "entriesCount": 1162,
              "volumeSizeBytes": 605343
            }
          ]
        },
        {
          "name": "KAFKA",
          "entriesCount": 8266,
          "volumeSizeBytes": 12059585,
          "color": "#000000",
          "methods": [
            {
              "name": "Produce",
              "entriesCount": 760,
              "volumeSizeBytes": 1304364
            },
            {
              "name": "Metadata",
              "entriesCount": 2706,
              "volumeSizeBytes": 3273481
            },
            {
              "name": "ApiVersions",
              "entriesCount": 3677,
              "volumeSizeBytes": 6398862
            },
            {
              "name": "CreateTopics",
              "entriesCount": 1123,
              "volumeSizeBytes": 1082878
            }
          ]
        },
        {
          "name": "gRPC",
          "entriesCount": 10829,
          "volumeSizeBytes": 48790409,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 10829,
              "volumeSizeBytes": 48790409
            }
          ]
        },
        {
          "name": "HTTP",
          "entriesCount": 19188,
          "volumeSizeBytes": 87259874,
          "color": "#244c5a",
          "methods": [
            {
              "name": "GET",
              "entriesCount": 19186,
              "volumeSizeBytes": 87253974
            },
            {
              "name": "POST",
              "entriesCount": 2,
              "volumeSizeBytes": 5900
            }
          ]
        },
        {
          "name": "GQL",
          "entriesCount": 24070,
          "volumeSizeBytes": 82144676,
          "color": "#e10098",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 24070,
              "volumeSizeBytes": 82144676
            }
          ]
        },
        {
          "name": "REDIS",
          "entriesCount": 14016,
          "volumeSizeBytes": 7577827,
          "color": "#a41e11",
          "methods": [
            {
              "name": "SET",
              "entriesCount": 4673,
              "volumeSizeBytes": 2542881
            },
            {
              "name": "GET",
              "entriesCount": 9343,
              "volumeSizeBytes": 5034946
            }
          ]
        }
      ],
      "timestamp": 1656532800000
    },
    {
      "protocols": [
        {
          "name": "AMQP",
          "entriesCount": 12765,
          "volumeSizeBytes": 9166687,
          "color": "#ff6600",
          "methods": [
            {
              "name": "queue bind",
              "entriesCount": 1116,
              "volumeSizeBytes": 581502
            },
            {
              "name": "connection close",
              "entriesCount": 1116,
              "volumeSizeBytes": 555837
            },
            {
              "name": "basic consume",
              "entriesCount": 1116,
              "volumeSizeBytes": 642852
            },
            {
              "name": "basic publish",
              "entriesCount": 2232,
              "volumeSizeBytes": 1765656
            },
            {
              "name": "connection start",
              "entriesCount": 1651,
              "volumeSizeBytes": 1667158
            },
            {
              "name": "basic deliver",
              "entriesCount": 3302,
              "volumeSizeBytes": 2706947
            },
            {
              "name": "queue declare",
              "entriesCount": 1116,
              "volumeSizeBytes": 606074
            },
            {
              "name": "exchange declare",
              "entriesCount": 1116,
              "volumeSizeBytes": 640661
            }
          ]
        },
        {
          "name": "KAFKA",
          "entriesCount": 8491,
          "volumeSizeBytes": 13333916,
          "color": "#000000",
          "methods": [
            {
              "name": "CreateTopics",
              "entriesCount": 1119,
              "volumeSizeBytes": 1080862
            },
            {
              "name": "ApiVersions",
              "entriesCount": 3836,
              "volumeSizeBytes": 6671764
            },
            {
              "name": "Produce",
              "entriesCount": 795,
              "volumeSizeBytes": 2271673
            },
            {
              "name": "Metadata",
              "entriesCount": 2741,
              "volumeSizeBytes": 3309617
            }
          ]
        },
        {
          "name": "HTTP",
          "entriesCount": 19192,
          "volumeSizeBytes": 87295897,
          "color": "#244c5a",
          "methods": [
            {
              "name": "GET",
              "entriesCount": 19192,
              "volumeSizeBytes": 87295897
            }
          ]
        },
        {
          "name": "gRPC",
          "entriesCount": 10657,
          "volumeSizeBytes": 48019503,
          "color": "#244c5a",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 10657,
              "volumeSizeBytes": 48019503
            }
          ]
        },
        {
          "name": "GQL",
          "entriesCount": 23814,
          "volumeSizeBytes": 81265370,
          "color": "#e10098",
          "methods": [
            {
              "name": "POST",
              "entriesCount": 23814,
              "volumeSizeBytes": 81265370
            }
          ]
        },
        {
          "name": "REDIS",
          "entriesCount": 14017,
          "volumeSizeBytes": 7578582,
          "color": "#a41e11",
          "methods": [
            {
              "name": "SET",
              "entriesCount": 4673,
              "volumeSizeBytes": 2542995
            },
            {
              "name": "GET",
              "entriesCount": 9344,
              "volumeSizeBytes": 5035587
            }
          ]
        }
      ],
      "timestamp": 1656525600000
    }
  ]
}


export const PROTOCOLS = ["ALL PROTOCOLS","gRPC", "REDIS", "HTTP", "GQL", "AMQP", "KFAKA"];
export const ALL_PROTOCOLS = PROTOCOLS[0];

export const TrafficStatsModal: React.FC<TrafficStatsModalProps> = ({ isOpen, onClose, getPieStatsDataApi, getTimelineStatsDataApi }) => {

  const modes = Object.keys(StatsMode).filter(x => !(parseInt(x) >= 0));
  const [statsMode, setStatsMode] = useState(modes[0]);
  const [selectedProtocol, setSelectedProtocol] = useState("ALL PROTOCOLS");
  const [pieStatsData, setPieStatsData] = useState(null);
  const [timelineStatsData, setTimelineStatsData] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const commonClasses = useCommonStyles();

  const getTrafficStats = useCallback(async () => {
    if (isOpen && getPieStatsDataApi) {
      (async () => {
        try {
          setIsLoading(true);
          // const pieData = await getPieStatsDataApi();
          setPieStatsData(dataMock.pie);
          // const timelineData = await getTimelineStatsDataApi();
          setTimelineStatsData(dataMock.timeline);
        } catch (e) {
          console.error(e)
        } finally {
          setIsLoading(false)
        }
      })()
    }
  }, [isOpen, getPieStatsDataApi, getTimelineStatsDataApi, setPieStatsData, setTimelineStatsData])

  useEffect(() => {
    getTrafficStats();
  }, [getTrafficStats])

  const refreshStats = debounce(() => {
    getTrafficStats();
  }, 500);

  return (
    <Modal
      aria-labelledby="transition-modal-title"
      aria-describedby="transition-modal-description"
      open={isOpen}
      onClose={onClose}
      closeAfterTransition
      BackdropComponent={Backdrop}
      BackdropProps={{ timeout: 500 }}>
      <Fade in={isOpen}>
        <Box sx={modalStyle}>
          <div className={styles.closeIcon}>
            <img src={closeIcon} alt="close" onClick={() => onClose()} style={{ cursor: "pointer", userSelect: "none" }} />
          </div>
          <div className={styles.headlineContainer}>
            <div className={styles.title}>Traffic Statistics</div>
            <Button style={{ marginLeft: "2%", textTransform: 'unset' }}
              startIcon={<img src={refreshIcon} className="custom" alt="refresh"></img>}
              size="medium"
              variant="contained"
              className={commonClasses.outlinedButton + " " + commonClasses.imagedButton}
              onClick={refreshStats}
            >
              Refresh
            </Button>
          </div>
          <div className={styles.mainContainer}>
            <div className={styles.selectContainer}>
              <div>
                <span style={{ marginRight: 15 }}>Breakdown By</span>
                <select className={styles.select} value={statsMode} onChange={(e) => setStatsMode(e.target.value)}>
                  {modes.map(mode => <option key={mode} value={mode}>{mode}</option>)}
                </select>
              </div>
              <div>
                <span style={{ marginRight: 15 }}>Protocol</span>
                <select className={styles.select} value={selectedProtocol} onChange={(e) => setSelectedProtocol(e.target.value)}>
                  {PROTOCOLS.map(protocol => <option key={protocol} value={protocol}>{protocol}</option>)}
                </select>
              </div>
            </div>
            <div>
              {isLoading ? <div style={{ textAlign: "center", marginTop: 20 }}>
                <img alt="spinner" src={spinnerImg} style={{ height: 50 }} />
              </div> :
                <div>
                  <TrafficPieChart pieChartMode={statsMode} data={pieStatsData} selectedProtocol={selectedProtocol}/>
                  <TimelineBarChart timeLineBarChartMode={statsMode} data={timelineStatsData} selectedProtocol={selectedProtocol}/>
                </div>}
            </div>
          </div>
        </Box>
      </Fade>
    </Modal>
  );
}
